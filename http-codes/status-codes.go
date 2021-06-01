// contains Tofa-specific http status codes replied to cloud services amid interaction
package http_codes

// When human allowed a REG or ASK request
const ACTION_ALLOWED = 270

// When human denied a REG or ASK request
const ACTION_REJECTED = 570

// When there is sync error between GUI and daemon.
// Should not but may occur amid oservapp.sRegisterService
// when GUI replies back to ws pipe
const CL_DA_CONFLICT = 571
