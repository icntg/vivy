package errno

const (
	NoError                          = 0
	HttpContinue                     = 100
	HttpSwitchingProtocols           = 101
	HttpOK                           = 200
	HttpCreated                      = 201
	HttpAccepted                     = 202
	HttpNonAuthoritativeInformation  = 203
	HttpNoContent                    = 204
	HttpResetContent                 = 205
	HttpPartialContent               = 206
	HttpMultipleChoices              = 300
	HttpMovedPermanently             = 301
	HttpFound                        = 302
	HttpSeeOther                     = 303
	HttpNotModified                  = 304
	HttpUseProxy                     = 305
	HttpUnused                       = 306
	HttpTemporaryRedirect            = 307
	HttpBadRequest                   = 400
	HttpUnauthorized                 = 401
	HttpPaymentRequired              = 402
	HttpForbidden                    = 403
	HttpNotFound                     = 404
	HttpMethodNotAllowed             = 405
	HttpNotAcceptable                = 406
	HttpProxyAuthenticationRequired  = 407
	HttpRequestTimeout               = 408
	HttpConflict                     = 409
	HttpGone                         = 410
	HttpLengthRequired               = 411
	HttpPreconditionFailed           = 412
	HttpRequestEntityTooLarge        = 413
	HttpRequestURITooLarge           = 414
	HttpUnsupportedMediaType         = 415
	HttpRequestedRangeNotSatisfiable = 416
	HttpExpectationFailed            = 417
	HttpInternalServerError          = 500
	HttpNotImplemented               = 501
	HttpBadGateway                   = 502
	HttpServiceUnavailable           = 503
	HttpGatewayTimeout               = 504
	HttpHTTPVersionNotSupported      = 505
)
