package handler

import (
	"WatchHive/pkg/helper"
	interfaces "WatchHive/pkg/helper/interface"
	service "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"

	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"errors"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

type AdminHandler struct {
	adminUseCase service.AdminUseCase
	helper       interfaces.Helper
}

func NewAdminHandler(usecase service.AdminUseCase, helper interfaces.Helper) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
		helper:       helper,
	}
}

// LoginHandler handles the login operation for an admin.
// @Summary Admin login
// @Description Authenticate an admin and get access token
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body models.AdminLogin true "Admin credentials for login"
// @Success 200 {object} response.Response "Admin login successful"
// @Failure 400 {object} response.Response "Invalid request or constraints not satisfied"
// @Failure 401 {object} response.Response "Unauthorized: cannot authenticate user"
// @Router /admin/ [post]
func (ad *AdminHandler) LoginHandler(c *gin.Context) {

	var adminDetails models.AdminLogin

	if err := c.BindJSON(&adminDetails); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgConstraintsErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	err := validator.New().Struct(adminDetails)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgConstraintsErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	admin, err := ad.adminUseCase.LoginHandler(adminDetails)
	if err != nil {
		errREsp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgAuthUserErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errREsp)
		return
	}

	c.Set("Access", admin.AccessToken)
	// c.Set("Refresh", admin.RefreshToken)

	successResp := response.ClientResponse(http.StatusOK, errmsg.MsgLoginSucces, admin, nil)
	c.JSON(http.StatusOK, successResp)

}

func (ad *AdminHandler) ValidateRefreshTokenAndCreateNewAccess(c *gin.Context) {
	refreshToken := c.Request.Header.Get("RefreshToken")

	// Check  refresh token is valid.

	_, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte("refreshsecret"), nil
	})
	if err != nil {
		c.AbortWithError(401, errors.New(errmsg.ErrRefreshToken))
		return
	}
	claims := &helper.AuthCustomClaims{
		Role: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	newAccessToken, err := token.SignedString([]byte("accesssecret"))
	if err != nil {
		c.AbortWithError(500, errors.New(errmsg.ErrAccessToken))
	}
	c.JSON(200, newAccessToken)
}

// BlockUser blocks a user by ID.
// @Summary Block a user
// @Description Blocks a user based on the provided ID
// @Tags Admin User Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id query string true "User ID to block"
// @Success 200 {object} response.Response "User blocked successfully"
// @Failure 400 {object} response.Response "Failed to block user"
// @Router /admin/users/block [patch]
func (ad *AdminHandler) BlockUser(c *gin.Context) {
	id := c.Query("id")
	err := ad.adminUseCase.BlockUser(id)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgUserBlockErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	succesResp := response.ClientResponse(http.StatusOK, errmsg.MsgUserBlockSucces, nil, nil)
	c.JSON(http.StatusOK, succesResp)

}

// UnBlockUser unblocks a user by ID.
// @Summary Unblock a user
// @Description Unblocks a user based on the provided ID
// @Tags Admin User Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id query string true "User ID to unblock"
// @Success 200 {object} response.Response "User unblocked successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to unblock user"
// @Router /admin/users/unblock [patch]
func (ad *AdminHandler) UnBlockUser(c *gin.Context) {

	id := c.Query("id")
	err := ad.adminUseCase.UnBlockUser(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgUserUnBlockErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgUserUnBlockSucces, nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetUsers retrieves users based on the provided page number.
// @Summary Retrieve users with pagination
// @Description Retrieves users based on the provided page number
// @Tags Admin User Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param page query int true "Page number for pagination"
// @Success 200 {object} response.Response "Users retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to retrieve users"
// @Router /admin/users [get]
func (ad *AdminHandler) GetUsers(c *gin.Context) {

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgPageNumFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	users, err := ad.adminUseCase.GetUsers(page)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGettingDataErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, users, nil)
	c.JSON(http.StatusOK, successRes)

}

// AdminDashBoard retrieves the dashboard information for admin.
// @Summary Retrieve admin dashboard information
// @Description Retrieves dashboard information for admin
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Success 200 {object} response.Response "Admin dashboard retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to retrieve dashboard"
// @Router /admin/dashboard [get]
func (ah *AdminHandler) AdminDashBoard(c *gin.Context) {
	dashboard, err := ah.adminUseCase.AdminDashboard()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGettingDataErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, dashboard, nil)
	c.JSON(http.StatusOK, successRes)
}

// FilteredSalesReport retrieves the  current sales report for a specified time period.
// @Summary Retrieve current sales report for a specific time period
// @Description Retrieves sales report for the specified time period
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param period query string true "Time period for sales report"
// @Success 200 {object} response.Response "Sales report retrieved successfully"
// @Failure 500 {object} response.Response "Unable to retrieve sales report"
// @Router /admin/currentsalesreport [get]
func (ah *AdminHandler) FilteredSalesReport(c *gin.Context) {

	timePeriod := c.Query("period")
	salesReport, err := ah.adminUseCase.FilteredSalesReport(timePeriod)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, errmsg.MsgGettingDataErr, nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	message := " current " + timePeriod + errmsg.MsgGetSucces

	success := response.ClientResponse(http.StatusOK, message, salesReport, nil)
	c.JSON(http.StatusOK, success)
}

//repot by date

// SalesReportByDate generates a sales report within a specified date range.
// @Summary Generate sales report by date range
// @Description Retrieves sales report data between the provided start and end dates.
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param start query string true "Start date (format: 'DD-MM-YYYY')" Format(DD-MM-YYYY)
// @Param end query string true "End date (format: 'DD-MM-YYYY')" Format(DD-MM-YYYY)
// @Failure 400 {object} response.Response "Bad request: Start or end date is empty"
// @Failure 500 {object} response.Response "Internal server error: Sales report could not be retrieved"
// @Success 200 {object} response.Response "Success: Sales report retrieved successfully"
// @Router /admin/salesreport [get]
func (ah *AdminHandler) SalesReportByDate(c *gin.Context) {
	startDateStr := c.Query("start")
	endDateStr := c.Query("end")
	if startDateStr == "" || endDateStr == "" {
		err := response.ClientResponse(http.StatusBadRequest, errmsg.MsgEmptyDateErr, nil, "Empty date string")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	report, err := ah.adminUseCase.ExecuteSalesReportByDate(startDateStr, endDateStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, errmsg.MsgGettingDataErr, nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	success := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, report, nil)
	c.JSON(http.StatusOK, success)
}

// SalesByDate gets sales details for a specific date and allows downloading the report in PDF or Excel format.
//
// @Summary Get sales details by date
// @Description Get sales details for a specific date and download the report in PDF or Excel format
// @Tags Admin Dashboard
// @security BearerTokenAuth
// @Param year query integer true "Year for sales data"
// @Param month query integer true "Month for sales data"
// @Param day query integer true "Day for sales data"
// @Param download query string false "Download format (pdf or excel)"
// @Success 200 {object} response.Response "Successfully retrieved sales details"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Failure 502 {object} response.Response "Bad Gateway"
// @Router /admin/printsales [get]
func (a *AdminHandler) PrintSalesByDate(c *gin.Context) {
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGetErr+"year", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	month := c.Query("month")
	monthInt, err := strconv.Atoi(month)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGetErr+"month", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	day := c.Query("day")
	dayInt, err := strconv.Atoi(day)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGetErr+"day", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	body, err := a.adminUseCase.SalesByDate(dayInt, monthInt, yearInt)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGettingDataErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	download := c.Query("download")
	if download == "pdf" {
		pdf, err := a.adminUseCase.PrintSalesReport(body)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, errmsg.MsgGettingDataErr, nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
		c.Header("Content-Disposition", "attachment;filename=totalsalesreport.pdf")

		pdfFilePath := "salesReport/totalsalesreport.pdf"

		err = pdf.OutputFileAndClose(pdfFilePath)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, errmsg.MsgPrintErr, nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}

		c.Header("Content-Disposition", "attachment; filename=total_sales_report.pdf")
		c.Header("Content-Type", "application/pdf")

		c.File(pdfFilePath)

		c.Header("Content-Type", "application/pdf")

		err = pdf.Output(c.Writer)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, errmsg.MsgPrintErr, nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
	} else {

		excel, err := a.helper.ConvertToExel(body)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, errmsg.MsgPrintErr, nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}

		fileName := "sales_report.xlsx"

		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		if err := excel.Write(c.Writer); err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, errmsg.MsgServErr, nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
	}

	succesRes := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, body, nil)
	c.JSON(http.StatusOK, succesRes)
}
