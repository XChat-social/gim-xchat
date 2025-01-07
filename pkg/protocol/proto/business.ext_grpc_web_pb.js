/**
 * @fileoverview gRPC-Web generated client stub for pb
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v5.29.0
// source: business.ext.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js')
const proto = {};
proto.pb = require('./business.ext_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.pb.BusinessExtClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.pb.BusinessExtPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.SignInReq,
 *   !proto.pb.SignInResp>}
 */
const methodDescriptor_BusinessExt_SignIn = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/SignIn',
  grpc.web.MethodType.UNARY,
  proto.pb.SignInReq,
  proto.pb.SignInResp,
  /**
   * @param {!proto.pb.SignInReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.SignInResp.deserializeBinary
);


/**
 * @param {!proto.pb.SignInReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.SignInResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.SignInResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.signIn =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/SignIn',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_SignIn,
      callback);
};


/**
 * @param {!proto.pb.SignInReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.SignInResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.signIn =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/SignIn',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_SignIn);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.GetUserReq,
 *   !proto.pb.GetUserResp>}
 */
const methodDescriptor_BusinessExt_GetUser = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/GetUser',
  grpc.web.MethodType.UNARY,
  proto.pb.GetUserReq,
  proto.pb.GetUserResp,
  /**
   * @param {!proto.pb.GetUserReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.GetUserResp.deserializeBinary
);


/**
 * @param {!proto.pb.GetUserReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.GetUserResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.GetUserResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.getUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/GetUser',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_GetUser,
      callback);
};


/**
 * @param {!proto.pb.GetUserReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.GetUserResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.getUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/GetUser',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_GetUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.UpdateUserReq,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_BusinessExt_UpdateUser = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/UpdateUser',
  grpc.web.MethodType.UNARY,
  proto.pb.UpdateUserReq,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.pb.UpdateUserReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.pb.UpdateUserReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.updateUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/UpdateUser',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_UpdateUser,
      callback);
};


/**
 * @param {!proto.pb.UpdateUserReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.updateUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/UpdateUser',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_UpdateUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.SearchUserReq,
 *   !proto.pb.SearchUserResp>}
 */
const methodDescriptor_BusinessExt_SearchUser = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/SearchUser',
  grpc.web.MethodType.UNARY,
  proto.pb.SearchUserReq,
  proto.pb.SearchUserResp,
  /**
   * @param {!proto.pb.SearchUserReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.SearchUserResp.deserializeBinary
);


/**
 * @param {!proto.pb.SearchUserReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.SearchUserResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.SearchUserResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.searchUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/SearchUser',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_SearchUser,
      callback);
};


/**
 * @param {!proto.pb.SearchUserReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.SearchUserResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.searchUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/SearchUser',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_SearchUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.pb.TwitterAuthorizeURLResp>}
 */
const methodDescriptor_BusinessExt_GetTwitterAuthorizeURL = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/GetTwitterAuthorizeURL',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.pb.TwitterAuthorizeURLResp,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.TwitterAuthorizeURLResp.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.TwitterAuthorizeURLResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.TwitterAuthorizeURLResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.getTwitterAuthorizeURL =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/GetTwitterAuthorizeURL',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_GetTwitterAuthorizeURL,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.TwitterAuthorizeURLResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.getTwitterAuthorizeURL =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/GetTwitterAuthorizeURL',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_GetTwitterAuthorizeURL);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.TwitterSignInReq,
 *   !proto.pb.TwitterSignInResp>}
 */
const methodDescriptor_BusinessExt_TwitterSignIn = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/TwitterSignIn',
  grpc.web.MethodType.UNARY,
  proto.pb.TwitterSignInReq,
  proto.pb.TwitterSignInResp,
  /**
   * @param {!proto.pb.TwitterSignInReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.TwitterSignInResp.deserializeBinary
);


/**
 * @param {!proto.pb.TwitterSignInReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.TwitterSignInResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.TwitterSignInResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.twitterSignIn =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/TwitterSignIn',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_TwitterSignIn,
      callback);
};


/**
 * @param {!proto.pb.TwitterSignInReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.TwitterSignInResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.twitterSignIn =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/TwitterSignIn',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_TwitterSignIn);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.pb.DailySignInResp>}
 */
const methodDescriptor_BusinessExt_DailySignIn = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/DailySignIn',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.pb.DailySignInResp,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.DailySignInResp.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.DailySignInResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.DailySignInResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.dailySignIn =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/DailySignIn',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_DailySignIn,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.DailySignInResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.dailySignIn =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/DailySignIn',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_DailySignIn);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.pb.TwitterFollowResp>}
 */
const methodDescriptor_BusinessExt_FollowTwitter = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/FollowTwitter',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.pb.TwitterFollowResp,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.TwitterFollowResp.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.TwitterFollowResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.TwitterFollowResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.followTwitter =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/FollowTwitter',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_FollowTwitter,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.TwitterFollowResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.followTwitter =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/FollowTwitter',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_FollowTwitter);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.GetTaskStatusReq,
 *   !proto.pb.GetTaskStatusResp>}
 */
const methodDescriptor_BusinessExt_GetTaskStatus = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/GetTaskStatus',
  grpc.web.MethodType.UNARY,
  proto.pb.GetTaskStatusReq,
  proto.pb.GetTaskStatusResp,
  /**
   * @param {!proto.pb.GetTaskStatusReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.GetTaskStatusResp.deserializeBinary
);


/**
 * @param {!proto.pb.GetTaskStatusReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.GetTaskStatusResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.GetTaskStatusResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.getTaskStatus =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/GetTaskStatus',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_GetTaskStatus,
      callback);
};


/**
 * @param {!proto.pb.GetTaskStatusReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.GetTaskStatusResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.getTaskStatus =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/GetTaskStatus',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_GetTaskStatus);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.ClaimTaskRewardReq,
 *   !proto.pb.ClaimTaskRewardResp>}
 */
const methodDescriptor_BusinessExt_ClaimTaskReward = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/ClaimTaskReward',
  grpc.web.MethodType.UNARY,
  proto.pb.ClaimTaskRewardReq,
  proto.pb.ClaimTaskRewardResp,
  /**
   * @param {!proto.pb.ClaimTaskRewardReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.ClaimTaskRewardResp.deserializeBinary
);


/**
 * @param {!proto.pb.ClaimTaskRewardReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.ClaimTaskRewardResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.ClaimTaskRewardResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.claimTaskReward =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/ClaimTaskReward',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_ClaimTaskReward,
      callback);
};


/**
 * @param {!proto.pb.ClaimTaskRewardReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.ClaimTaskRewardResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.claimTaskReward =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/ClaimTaskReward',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_ClaimTaskReward);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.FillInviteCodeReq,
 *   !proto.pb.FillInviteCodeResp>}
 */
const methodDescriptor_BusinessExt_FillInviteCode = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/FillInviteCode',
  grpc.web.MethodType.UNARY,
  proto.pb.FillInviteCodeReq,
  proto.pb.FillInviteCodeResp,
  /**
   * @param {!proto.pb.FillInviteCodeReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.FillInviteCodeResp.deserializeBinary
);


/**
 * @param {!proto.pb.FillInviteCodeReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.FillInviteCodeResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.FillInviteCodeResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.fillInviteCode =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/FillInviteCode',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_FillInviteCode,
      callback);
};


/**
 * @param {!proto.pb.FillInviteCodeReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.FillInviteCodeResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.fillInviteCode =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/FillInviteCode',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_FillInviteCode);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.WalletSignInReq,
 *   !proto.pb.WalletSignInResp>}
 */
const methodDescriptor_BusinessExt_WalletSignIn = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/WalletSignIn',
  grpc.web.MethodType.UNARY,
  proto.pb.WalletSignInReq,
  proto.pb.WalletSignInResp,
  /**
   * @param {!proto.pb.WalletSignInReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.WalletSignInResp.deserializeBinary
);


/**
 * @param {!proto.pb.WalletSignInReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.WalletSignInResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.WalletSignInResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.walletSignIn =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/WalletSignIn',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_WalletSignIn,
      callback);
};


/**
 * @param {!proto.pb.WalletSignInReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.WalletSignInResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.walletSignIn =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/WalletSignIn',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_WalletSignIn);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pb.ModifyTaskStatusReq,
 *   !proto.pb.ModifyTaskStatusResp>}
 */
const methodDescriptor_BusinessExt_ModifyTaskStatus = new grpc.web.MethodDescriptor(
  '/pb.BusinessExt/ModifyTaskStatus',
  grpc.web.MethodType.UNARY,
  proto.pb.ModifyTaskStatusReq,
  proto.pb.ModifyTaskStatusResp,
  /**
   * @param {!proto.pb.ModifyTaskStatusReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pb.ModifyTaskStatusResp.deserializeBinary
);


/**
 * @param {!proto.pb.ModifyTaskStatusReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pb.ModifyTaskStatusResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pb.ModifyTaskStatusResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pb.BusinessExtClient.prototype.modifyTaskStatus =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pb.BusinessExt/ModifyTaskStatus',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_ModifyTaskStatus,
      callback);
};


/**
 * @param {!proto.pb.ModifyTaskStatusReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pb.ModifyTaskStatusResp>}
 *     Promise that resolves to the response
 */
proto.pb.BusinessExtPromiseClient.prototype.modifyTaskStatus =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pb.BusinessExt/ModifyTaskStatus',
      request,
      metadata || {},
      methodDescriptor_BusinessExt_ModifyTaskStatus);
};


module.exports = proto.pb;

