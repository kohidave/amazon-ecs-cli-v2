// Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
'use strict';

const aws = require('aws-sdk');

// priorityForRootRule is the max priority number that's always set for the listener rule that matches the root path "/"
const priorityForRootRule = "50000"

// These are used for test purposes only
let defaultResponseURL;

/**
 * Upload a CloudFormation response object to S3.
 *
 * @param {object} event the Lambda event payload received by the handler function
 * @param {object} context the Lambda context received by the handler function
 * @param {string} responseStatus the response status, either 'SUCCESS' or 'FAILED'
 * @param {string} physicalResourceId CloudFormation physical resource ID
 * @param {object} [responseData] arbitrary response data object
 * @param {string} [reason] reason for failure, if any, to convey to the user
 * @returns {Promise} Promise that is resolved on success, or rejected on connection error or HTTP error response
 */
let report = function (event, context, responseStatus, physicalResourceId, responseData, reason) {
    return new Promise((resolve, reject) => {
        const https = require('https');
        const {
            URL
        } = require('url');

        var responseBody = JSON.stringify({
            Status: responseStatus,
            Reason: reason,
            PhysicalResourceId: physicalResourceId || context.logStreamName,
            StackId: event.StackId,
            RequestId: event.RequestId,
            LogicalResourceId: event.LogicalResourceId,
            Data: responseData
        });

        const parsedUrl = new URL(event.ResponseURL || defaultResponseURL);
        const options = {
            hostname: parsedUrl.hostname,
            port: 443,
            path: parsedUrl.pathname + parsedUrl.search,
            method: 'PUT',
            headers: {
                'Content-Type': '',
                'Content-Length': responseBody.length
            }
        };

        https.request(options)
            .on('error', reject)
            .on('response', res => {
                res.resume();
                if (res.statusCode >= 400) {
                    reject(new Error(`Error ${res.statusCode}: ${res.statusMessage}`));
                } else {
                    resolve();
                }
            })
            .end(responseBody, 'utf8');
    });
};


const describeEnvironmentStack = async function (cf, stackName) {
    const stacks = await cf.describeStacks({
        StackName: stackName
    }).promise();
    return stacks.Stacks[0];
}
/**
 * Lists all the existing rules for a ALB Listener, finds the max of their
 * priorities, and then returns max + 1.
 *
 * @param {string} listenerArn the ARN of the ALB listener.

 * @returns {number} The next available ALB listener rule priority.
 */
const addEnvFeatures = async function (properties) {
    var cloudformation = new aws.CloudFormation();
    const stackName =  `${properties.Project}-${properties.Env}`;
    const envStack = await describeEnvironmentStack(cloudformation, stackName);
    envStack.Parameters.forEach(param => {
        const newFeatureValue = properties[param.ParameterKey];
        if(newFeatureValue !== undefined) {
            param.ParameterValue = newFeatureValue;
        }
    })

    try {
        await cloudformation.updateStack({
            UsePreviousTemplate: true,
            StackName: stackName,
            Parameters: envStack.Parameters,
            Capabilities: envStack.Capabilities,
            Tags: envStack.Tags,
        }).promise();
    } catch (err) {
        if(err.message !== 'No updates are to be performed.') {
            throw err;
        }
    }
    await cloudformation.waitFor('stackUpdateComplete', {StackName: stackName}).promise();
    const updatedEnvStack = await describeEnvironmentStack(cloudformation, stackName);
    const exportedValues = {};
    updatedEnvStack.Outputs.forEach(output => {
        var exportName = output.ExportName.replace(`${stackName}-`, '');
        exportedValues[exportName] = output.OutputValue;
        //exportedValues[exportName] = output.ExportName;

    })
    return exportedValues;
}

/**
 * Next Available ALB Rule Priority handler, invoked by Lambda
 */
exports.addAppFeaturesToEnv = async function(event, context) {
    var responseData = {};
    var physicalResourceId = `${event.ResourceProperties.Project}/${event.ResourceProperties.Env}/${event.ResourceProperties.App}/features`;
    var featuresSummary = "";
    try {
      switch (event.RequestType) {
        case 'Update':
        case 'Delete':
        case 'Create':
          responseData = await addEnvFeatures(event.ResourceProperties);
          break;
        default:
          throw new Error(`Unsupported request type ${event.RequestType}`);
      }
      await report(event, context, 'SUCCESS', physicalResourceId, responseData, featuresSummary);
    } catch (err) {
      console.log(`Caught error ${err}.`);
      await report(event, context, 'FAILED', physicalResourceId, null, err.message);
    }
  };


/**
 * @private
 */
exports.withDefaultResponseURL = function(url) {
  defaultResponseURL = url;
};
