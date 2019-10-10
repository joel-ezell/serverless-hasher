This service enables hashing of numbers and retrieving those hashes. It also provides a stats endpoint for retrieving statistics about those operations.

The code doesn't have any concurrency protection because that's all handled by DynamoDB. There's also no need to create a new GoRoutine for every long-running operation because the Lambda infrastructure handles all of that.

The HTTP front-end is facilitated by Amazon API Gateway, which has been configured to invoke the lambda function upon receipt of API messages to the various endpoints.

This particular architecture is completely elastic. For all intents and purposes, it scales up nearly infinitely. It will automatically scale down to 0 when there is no traffic, though there is a "cold-start" delay incurred on first invocation to a Lambda that has scaled down completely.

The service would likely be very inexpensive to operate, particularly if the operation took a shorter time than 5 seconds to complete. Serverless is especially attractive if there is relatively light load, or if the load is highly variable.
