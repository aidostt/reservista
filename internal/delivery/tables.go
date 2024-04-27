package delivery

import (
	proto_table "github.com/aidostt/protos/gen/go/reservista/table"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (h *Handler) table(api *gin.RouterGroup) {
	tables := api.Group("tables")
	{
		tables.GET("/view", h.getTable)
		tables.GET("/all/restaurant", h.getTablesByRestId)
		tables.POST("/add", h.addTable)
		tables.DELETE("/delete", h.deleteTableById)
		tables.GET("/all/restaurant/available", h.getAvailableTables)
		tables.GET("/all/restaurant/reserved", h.getReservedTables)
		tables.PATCH("/update", h.updateTableById)
	}
}

func (h *Handler) getTable(c *gin.Context) {
	var input idInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	conn, err := h.Dialog.NewConnection(h.Dialog.Addresses.Reservations)
	defer conn.Close()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong...")
		return
	}
	client := proto_table.NewTableClient(conn)

	table, err := client.GetTable(c.Request.Context(), &proto_table.IDRequest{Id: input.Id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a gRPC status error
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
			return
		}
		switch st.Code() {
		case codes.Internal:
			newResponse(c, http.StatusInternalServerError, "microservice failed to execute functionality:"+err.Error())
		default:
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *Handler) getTablesByRestId(c *gin.Context) {
	var input idInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	conn, err := h.Dialog.NewConnection(h.Dialog.Addresses.Reservations)
	defer conn.Close()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong...")
		return
	}
	client := proto_table.NewTableClient(conn)

	tables, err := client.GetTablesByRestId(c.Request.Context(), &proto_table.IDRequest{Id: input.Id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a gRPC status error
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			newResponse(c, http.StatusBadRequest, "invalid argument: "+err.Error())
		case codes.Internal:
			newResponse(c, http.StatusInternalServerError, "microservice failed to execute functionality:"+err.Error())
		default:
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, tables)
}

func (h *Handler) addTable(c *gin.Context) {
	var input tableInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	conn, err := h.Dialog.NewConnection(h.Dialog.Addresses.Reservations)
	defer conn.Close()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong...")
		return
	}
	client := proto_table.NewTableClient(conn)

	statusResponse, err := client.AddTable(c.Request.Context(), &proto_table.AddTableRequest{
		NumberOfSeats: input.NumberOfSeats,
		TableNumber:   input.TableNumber,
		RestaurantID:  input.RestaurantID,
		IsReserved:    input.IsReserved,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a gRPC status error
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
			return
		}
		switch st.Code() {
		case codes.Internal:
			newResponse(c, http.StatusInternalServerError, "microservice failed to execute functionality:"+err.Error())
		default:
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
		}
		return
	}
	if !statusResponse.GetStatus() {
		newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": statusResponse.Status})
}

func (h *Handler) deleteTableById(c *gin.Context) {
	var input idInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	conn, err := h.Dialog.NewConnection(h.Dialog.Addresses.Reservations)
	defer conn.Close()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong...")
		return
	}
	client := proto_table.NewTableClient(conn)

	table, err := client.DeleteTableById(c.Request.Context(), &proto_table.IDRequest{Id: input.Id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a gRPC status error
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
			return
		}
		switch st.Code() {
		case codes.Internal:
			newResponse(c, http.StatusInternalServerError, "microservice failed to execute functionality:"+err.Error())
		default:
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *Handler) getAvailableTables(c *gin.Context) {
	var input idInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	conn, err := h.Dialog.NewConnection(h.Dialog.Addresses.Reservations)
	defer conn.Close()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong...")
		return
	}
	client := proto_table.NewTableClient(conn)

	table, err := client.GetAvailableTables(c.Request.Context(), &proto_table.IDRequest{Id: input.Id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a gRPC status error
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
			return
		}
		switch st.Code() {
		case codes.Internal:
			newResponse(c, http.StatusInternalServerError, "microservice failed to execute functionality:"+err.Error())
		default:
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *Handler) getReservedTables(c *gin.Context) {
	var input idInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	conn, err := h.Dialog.NewConnection(h.Dialog.Addresses.Reservations)
	defer conn.Close()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong...")
		return
	}
	client := proto_table.NewTableClient(conn)

	table, err := client.GetReservedTables(c.Request.Context(), &proto_table.IDRequest{Id: input.Id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a gRPC status error
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
			return
		}
		switch st.Code() {
		case codes.Internal:
			newResponse(c, http.StatusInternalServerError, "microservice failed to execute functionality:"+err.Error())
		default:
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *Handler) updateTableById(c *gin.Context) {
	var input tableInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	conn, err := h.Dialog.NewConnection(h.Dialog.Addresses.Reservations)
	defer conn.Close()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong...")
		return
	}
	client := proto_table.NewTableClient(conn)

	statusResponse, err := client.UpdateTableById(c.Request.Context(), &proto_table.UpdateTableRequest{
		Id:            input.Id,
		NumberOfSeats: input.NumberOfSeats,
		IsReserved:    input.IsReserved,
		TableNumber:   input.TableNumber,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a gRPC status error
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
			return
		}
		switch st.Code() {
		case codes.Internal:
			newResponse(c, http.StatusInternalServerError, "microservice failed to execute functionality:"+err.Error())
		default:
			newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
		}
		return
	}
	if !statusResponse.GetStatus() {
		newResponse(c, http.StatusInternalServerError, "unknown error when calling sign up:"+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": statusResponse.Status})
}
