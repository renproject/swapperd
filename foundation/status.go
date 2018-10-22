package foundation

type Status int64

var INACTIVE = Status(0)
var INITIATED = Status(1)
var AUDITED = Status(2)
var AUDIT_FAILED = Status(3)
var REDEEMED = Status(4)
var REFUNDED = Status(5)
