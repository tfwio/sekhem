// session package works with GORM to enable `http.Cookie`s,
// `Session` and `User` capability.
//
// This includes some cryptography that allows for a user's password
// and verification.
//
// Currently, this system is newly written and limits a User to utilizing
// one session where the Session is re-used and the user profile is locked
// into using a particular IP (client-machine).
//
// Theoretically, there is some advantge to locking down "security" in
// this manner.
//
// Changes will be documented here in the future when sessions privelages
// change.
//

package session
