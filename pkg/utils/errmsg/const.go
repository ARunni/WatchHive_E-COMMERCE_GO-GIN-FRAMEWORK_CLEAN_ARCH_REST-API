package errmsg

// Error strings

const ErrRefreshToken = "refresh token is sinvalid  login again "
const ErrGetDB = "error in retriving data from database"
const ErrWriteDB = "error in writing data to database"
const ErrUpdateDB = "error in updating data to database"
const ErrAccessToken = "error in creating new accesstoken"
const ErrServer = "internal server error"
const ErrOtpValidate = "failed to validate otp"
const ErrDataIsNot = "data is not"
const ErrAlphabet = "data contains non-alphabetic characters"
const ErrUserExistFalse = "user does not exist"
const ErrDb = "data base error"
const ErrProductExist = "product does not exist"
const ErrUserExist = "user does not exist"
const ErrProductExistTrue = "product already exists in the database"
const ErrIdExist = "id does not exist"
const ErrCartFalse = "no cart found"
const ErrDbConnect = "database connection is nil"
const ErrDatatypeConversion = "convertion of datatype is failed"
const ErrOfferExistTrue = "offer already exist"
const ErrGetOffer = "error in getting offers of"
const ErrGetData = "error in getting data"
const ErrDataNegative = "data cannot be negative"
const ErrBlockAlready = "already blocked"
const ErrUnBlockAlready = "already unblocked"
const ErrFieldEmpty = "field cannot be empty"
const ErrInvalidTimePeriod = "invalid time period, available options : week, month & year"
const ErrFormat = "enter the data in correct format"
const ErrInvalidFormat = "invalid format"
const ErrInvalidPId = "invalid product id"
const ErrInvalidOId = "invalid order id"
const ErrDataZero = "data must be 1 or greater"
const ErrOutOfStock = "out of stock"
const ErrLimitExceeds = "limit exceeds"
const ErrEmptyCart = "cart is empty"
const ErrCartProductExist = "product not available in cart"
const ErrExistTrue = "already exist"
const ErrInvalidCId = "invalid category id"
const ErrCatExistFalse = "category does not exist"
const ErrInvalidData = "invalid data"
const ErrOfferAdd = "error in adding offer"
const ErrNotExist = "does not exist"
const ErrUserOwnedOrder = "the order is not done by this user"
const ErrCancelAlready = "the order is already cancelled, so no point in cancelling"
const ErrReturnedAlready = "the order is already returned"
const ErrCancelAlreadyReturn = "the order is cancelled,cannot return it"
const ErrDeliveredAlreadyCancel = "the order is delivered cannot be cancelled"
const ErrDeliveredAlready = "the order is delivered, you can return it"
const ErrCancelAlreadyApprove = "the order is cancelled,cannot approve it"
const ErrPendingApprove = "the order is pending,cannot approve it"
const ErrDeliveredApprove = "the order is already deliverd"
const ErrPendingReturn = "the order is pending,cannot return it"
const ErrProcessingReturn = "the order is processing cannot return it"
const ErrShippedReturn = "the order is shipped cannot return it"
const ErrDeliverInvoice = "wait for the invoice until the product is received"
const ErrInvalidPhone = "invalid phone number"
const ErrVerify = "error while verifying"
const ErrAlreadyPaid = "already paid"
const ErrAlreadyUser = "user already exist, sign in"
const ErrPasswordMatch = "password does not match"
const ErrPasswordHash = "error hashing password"
const ErrCreateRefferal = "referral creation failed"
const ErrCreateTocken = "could not create token"
const ErrInternal = "internal error"
const ErrUserBlockTrue = "user is blocked"
const ErrPassword = "password is incorrect"
const ErrInvalidName = "invalid name"
const ErrInvalidPin = "invalid pin number"
const ErrInvalidUId = "invalid user id"
const ErrChangePassword = "password cannot change"


const StatusApprove = "approved"

// Message strings
const DbErr = "Data base error"

const MsgConstraintsErr = "Constraints not satisfied"
const MsgAuthUserErr = "Cannot authenticate user"
const MsgFormatErr = "Details is not in correct format"
const MsgLoginSucces = "Logined successfully"
const MsgUserBlockErr = "User could not be blocked"
const MsgUserBlockSucces = "Successfully blocked the user"
const MsgUserUnBlockErr = "User could not be unblocked"
const MsgUserUnBlockSucces = "Successfully unblocked the user"
const MsgPageNumFormatErr = "page number not in right format"
const MsgPageCountErr = "page count not in right format"
const MsgGettingDataErr = "could not retrieve Data"
const MsgGetSucces = "Successfully retrieved the Data"
const MsgEmptyDateErr = "Start or End date is empty"
const MsgGetErr = "error in getting"
const MsgPrintErr = "error in printing"
const MsgServErr = "Error in serving the sales report"
const MsgSuccess = "Success"
const MsgBadRequestErr = "bad request"
const MsgAddSuccess = "Successfully Added"
const MsgAddCartErr = "Cannot Add to Cart"
const MsgListErr = "Cannot list data"
const MsgListingErr = "Product cannot be displayed"
const MsgUpdateQuantityErr = "Cannot update quantity"
const MsgQuantityUpdationFailErr = "Updating quantity Failed"
const MsgRemoveCartErr = "Removing from cart is Failed"
const MsgAddErr = "Could not add"
const MsgUpdateErr = "Could not update"
const MsgCouponExpiryErr = "Coupon cannot be made invalid"
const MsgOTPSentErr = "OTP not sent"
const MsgOTPVerifyErr = "Could not verify OTP"
const MsgOTPSentSuccess = "OTP sent successfully"
const MsgOTPVerifySuccess = "Successfully verified OTP"
const MsgPaymentErr = "Cannot make payment"
const MsgErr = "error"
const MsgUpdateSuccess = "Successfully Updated"
const MsgUserIdErr = "user_id not found"
const MsgRequiredUserIdErr = "user_id is required"
const MsgInvalidIdErr = "invalid user_id type"
const MsgIdDatatypeErr = "user_id must be an integer"
const MsgEditErr = "could not edit the data"
const MsgStockUpdateErr = "Could  not update the product stock"
const MsgConvErr = "conversion error"
const MsdGetIdErr = "Failed to get user id"
const MsgLoginErr = "user could not be logged in"
const MsgSignUperr = "user could not signed up"
const MsgIdErr = "error in reading the order id"
const MsgCanelErr = "Couldn't cancel the order"
const MsgOrderApproveErr = "Couldn't approve the order"
const MsgOrderErr = "Could not do the order"
const MsgCheckoutErr = "CheckOut Failed"
const MsgTokenErr = "Invalid Authorization Token error"
const MsgTokenMissingErr = "Missing authorization token"
const MsgUnAuthErr = "Unauthorized access"
const MsgIdGetErr = "Error retrieving ID"
