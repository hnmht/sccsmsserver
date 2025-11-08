package handlers

import (
	"encoding/base64"
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"
	"sccsmsserver/pkg/minio"
	"sccsmsserver/pkg/security"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// User login handler
func LoginHandler(c *gin.Context) {
	// Step 1: Get request parameters
	p := new(pg.ParamLogin)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Provde additional login information
	p.ClientIP = c.ClientIP()
	p.ClientType = c.Request.Header.Get("XClientType")
	p.UserAgent = c.Request.UserAgent()

	// User login validation
	resStatus, token, _ := pg.Login(p)
	// Respond to client request
	ResponseWithMsg(c, resStatus, token)
}

// Get User Information handler
func UserInfoHandler(c *gin.Context) {
	userID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, nil)
		return
	}
	var u = pg.User{ID: userID}
	resStatus, _ = u.GetUserInfoByID()
	ResponseWithMsg(c, resStatus, u)
}

// Add user handler
func AddUserHandler(c *gin.Context) {
	u := new(pg.User)
	err := c.ShouldBind(u)
	if err != nil {
		zap.L().Error("AddUserHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator id
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, u)
		return
	}
	u.Creator.ID = operatorID

	// Decrypt the RSA password sent from the frent end
	op, _ := base64.StdEncoding.DecodeString(u.Password)
	oriPassword, err := security.ScRsa.Decrypt(op)
	if err != nil {
		zap.L().Error("Decrypt password failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, nil)
		return
	}
	u.Password = string(oriPassword)
	// Add user
	resStatus, _ = u.Add()
	// Response
	ResponseWithMsg(c, resStatus, u)
}

// Edit user handler
func EditUserHandler(c *gin.Context) {
	u := new(pg.User)
	err := c.ShouldBind(u)
	if err != nil {
		zap.L().Error("EditUserHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// If the password field is not empty,
	// it indicates that the user's password needs to be changed.
	if u.Password != "" {
		// Decrypt the RSA password sent from the frent end.
		op, _ := base64.StdEncoding.DecodeString(u.Password)
		oriPassword, err := security.ScRsa.Decrypt(op)
		if err != nil {
			zap.L().Error("Decrypt password failed", zap.Error(err))
			ResponseWithMsg(c, i18n.CodeInternalError, nil)
			return
		}
		u.Password = string(oriPassword)
	}
	// Get operator id
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, u)
		return
	}
	u.Modifier.ID = operatorID
	// Edit user
	resStatus, _ = u.Edit()
	// Resopnse
	ResponseWithMsg(c, resStatus, u)
}

// User Updates VIA personal center handler
func ModifyProfileHandler(c *gin.Context) {
	u := new(pg.User)
	err := c.ShouldBind(u)
	if err != nil {
		zap.L().Error("ModifyProfileHandler  invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get opertor id
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, u)
		return
	}
	u.Modifier.ID = operatorID
	// Update
	resStatus, _ = u.ModifyProfile()
	// Resoponse
	ResponseWithMsg(c, resStatus, u)
}

// Logout handler
func LogoutHandler(c *gin.Context) {
	// Get Client type
	clientType := c.Request.Header.Get("XClientType")
	// Get operator id
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, nil)
		return
	}
	// Delete online user from local cache
	var ou pg.OnlineUser
	ou.ClientType = clientType
	ou.User.ID = operatorID
	_, err := ou.Del()
	if err != nil {
		zap.L().Error("LogoutHandler ou.Del failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, err)
		return
	}
	ResponseWithMsg(c, i18n.StatusOK, nil)
}

// Get User list handler
func GetUsersHandler(c *gin.Context) {
	users, resStatus, _ := pg.GetUsers()
	ResponseWithMsg(c, resStatus, users)
}

// Delete user handler
func DeleteUserHandler(c *gin.Context) {
	u := new(pg.User)
	err := c.ShouldBind(u)
	if err != nil {
		zap.L().Error("DeleteUserHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator id
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, u)
		return
	}
	u.Modifier.ID = operatorID
	// Delete user
	resStatus, _ = u.Delete()
	ResponseWithMsg(c, resStatus, u)
}

// Batch delete user handler
func DeleteUsersHandler(c *gin.Context) {
	users := new([]pg.User)
	err := c.ShouldBind(users)
	if err != nil {
		zap.L().Error("DeleteUsersHandler invalid params", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator id
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, users)
		return
	}
	// Batch Delete
	resStatus, _ = pg.DeleteUsers(users, operatorID)
	// Response
	ResponseWithMsg(c, resStatus, users)
}

// Check user name exists handler
func CheckUserNameExistHandler(c *gin.Context) {
	u := new(pg.User)
	err := c.ShouldBind(u)
	if err != nil {
		zap.L().Error("CheckUserNameExistHandler  invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	resStatus, _ := u.CheckUserNameExist()
	ResponseWithMsg(c, resStatus, u)
}

// Check user code exists handler
func CheckUserCodeExistHandler(c *gin.Context) {
	u := new(pg.User)
	err := c.ShouldBind(u)
	if err != nil {
		zap.L().Error("CheckUserCodeExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check
	resStatus, _ := u.CheckUserCodeExist()
	// Resopnse
	ResponseWithMsg(c, resStatus, u)
}

// Change user password handler
func ChangeUserPasswordHandler(c *gin.Context) {
	p := new(pg.ParamChangePwd)
	err := c.ShouldBind(p)
	if err != nil {
		zap.L().Error("ChangeUserPasswordHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Decrypt the RSA field sent from the frent end
	op, err := base64.StdEncoding.DecodeString(p.Password)
	if err != nil {
		zap.L().Error("ChangeUserPasswordHandler base64.StdEncoding.DecodeString(p.Password) failed: ", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, err)
		return
	}
	np, err := base64.StdEncoding.DecodeString(p.NewPassword)
	if err != nil {
		zap.L().Error("ChangeUserPasswordHandler base64.StdEncoding.DecodeString(p.NewPassword) failed: ", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, err)
		return
	}
	cnp, err := base64.StdEncoding.DecodeString(p.ConfirmNewPwd)
	if err != nil {
		zap.L().Error("ChangeUserPasswordHandler base64.StdEncoding.DecodeString(p.ConfirmNewPwd) failed: ", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, err)
		return
	}
	oriPwd, err := security.ScRsa.Decrypt(op)
	if err != nil {
		zap.L().Error("ChangeUserPasswordHandler security.ScRsa.Decrypt(op) failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, err)
		return
	}
	oriNewPwd, err := security.ScRsa.Decrypt(np)
	if err != nil {
		zap.L().Error("ChangeUserPasswordHandler security.ScRsa.Decrypt(np) failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, err)
		return
	}
	oriConfirmNewPwd, err := security.ScRsa.Decrypt(cnp)
	if err != nil {
		zap.L().Error("ChangeUserPasswordHandler security.ScRsa.Decrypt(cnp) failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, err)
		return
	}

	p.Password = string(oriPwd)
	p.NewPassword = string(oriNewPwd)
	p.ConfirmNewPwd = string(oriConfirmNewPwd)

	resStatus, _ := p.ChangePassword()

	ResponseWithMsg(c, resStatus, nil)
}

// Change User avatar handler
func ChangeUserAvatarHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("ChangeUserAvatarHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	fileName := file.Filename
	fileObj, err := file.Open()
	if err != nil {
		zap.L().Error("ChangeUserAvatarHandler file.Open failed:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Upload the file to the MINIO server
	_, err = minio.UploadFile(fileName, fileObj, file.Size)
	if err != nil {
		zap.L().Error("ChangeUserAvatarHandler minio.UploadFile failed:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get the file URL to the MINIO server
	presignedURL, err := minio.GetFileUrl(fileName, time.Second*24*60*60)
	if err != nil {
		zap.L().Error("ChangeUserAvatarHandler  minio.GetFileUrl failed:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	ResponseWithMsg(c, i18n.StatusOK, presignedURL)
}
