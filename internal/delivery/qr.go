package delivery

import (
	proto_qr "github.com/aidostt/protos/gen/go/reservista/qr"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (h *Handler) qr(api *gin.RouterGroup) {
	qr := api.Group("/qr")
	{
		authenticated := qr.Group("/", h.userIdentity)
		authenticated.Use(h.isExpired)
		{
			authenticated.POST("/generate", h.generateQR)
			//TODO: id = reservation ID
			authenticated.GET("/scan/:id", h.scanQR)
		}
	}
}

func (h *Handler) generateQR(c *gin.Context) {
	var inp qrInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	conn, err := h.Dialog.NewConnection(h.Dialog.Addresses.QRs)
	defer conn.Close()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong...")
		return
	}
	//TODO: validate qr content
	qr := proto_qr.NewQRClient(conn)
	resp, err := qr.Generate(c.Request.Context(), &proto_qr.GenerateRequest{
		Content: "http://" + h.HttpAddress + "/api/qr/scan/" + inp.ReservationID,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a gRPC status error
			newResponse(c, http.StatusInternalServerError, "unknown error when calling generate qr:"+err.Error())
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			newResponse(c, http.StatusBadRequest, "invalid argument")
		case codes.Internal:
			newResponse(c, http.StatusInternalServerError, "microservice failed to execute functionality:"+err.Error())
		default:
			newResponse(c, http.StatusInternalServerError, "unknown error when calling generate qr:"+err.Error())
		}
		return
	}
	c.Header("Content-Type", "image/png")
	c.Writer.Write(resp.GetQR())
}

func (h *Handler) scanQR(c *gin.Context) {
	reservationID := c.Param("id")
	if reservationID == "" {
		newResponse(c, http.StatusBadRequest, "missing ID in the URL")
		return
	}
	conn, err := h.Dialog.NewConnection(h.Dialog.Addresses.QRs)
	defer conn.Close()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong...")
		return
	}
	qr := proto_qr.NewQRClient(conn)
	//TODO: check whether user from context is admin/restaurant stuff
	userID, ok := c.Get(userCtx)
	if !ok {
		newResponse(c, http.StatusBadRequest, "unauthorized access")
		return
	}
	resp, err := qr.Scan(c.Request.Context(), &proto_qr.ScanRequest{
		UserID:        userID.(string),
		ReservationID: reservationID,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a gRPC status error
			newResponse(c, http.StatusInternalServerError, "unknown error when calling scan qr:"+err.Error())
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			newResponse(c, http.StatusBadRequest, "invalid argument")
		case codes.Internal:
			newResponse(c, http.StatusInternalServerError, "microservice failed to execute functionality:"+err.Error())
		case codes.Unauthenticated:
			newResponse(c, http.StatusInternalServerError, "unauthorized access")
		default:
			newResponse(c, http.StatusInternalServerError, "unknown error when calling scan qr:"+err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}
