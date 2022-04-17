package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/murtaza-udaipurwala/pseudocoin/core"
	"github.com/murtaza-udaipurwala/pseudocoin/jsonrpc"
)

type IService interface {
	CreateWallet() (*Wallet, error)
	GetBalance(string) (*jsonrpc.Balance, error)
	Send(*Send, string) (*jsonrpc.Send, error)
	GetBlocks(*BlockQuery) (*jsonrpc.Blocks, error)
	GetAddress(string) (string, error)
	GetMyTXs(addr string) (*jsonrpc.MyTXs, error)
}

type Controller struct {
	s IService
}

func NewController(s IService) *Controller {
	return &Controller{s}
}

func (c *Controller) CreateWallet(ctx *fiber.Ctx) error {
	w, err := c.s.CreateWallet()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"successful": false,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"public_key":  w.PubKey,
		"private_key": w.PrivKey,
		"successful":  true,
	})
}

func (c *Controller) GetBalance(ctx *fiber.Ctx) error {
	addr := ctx.Query("addr")
	if !core.ValidateAddress(addr) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error":      "invalid address",
		})
	}

	bal, err := c.s.GetBalance(addr)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"successful": false,
			"error":      err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"successful": true,
		"balance":    bal.Balance,
		"address":    bal.Address,
	})
}

func (c *Controller) Send(ctx *fiber.Ctx) error {
	var req Send
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error":      err.Error(),
		})
	}

	if !core.ValidateAddress(req.RecvAddr) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error":      "invalid receiver's address",
		})
	}

	w := core.Wallet{}
	err = w.DecodePubKeys(req.SenderPub)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error":      err.Error(),
		})
	}

	sender, err := w.GetAddress()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error":      err.Error(),
		})
	}

	res, err := c.s.Send(&req, sender)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"successful": false,
			"error":      err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"successful": true,
		"msg":        res.Msg,
	})
}

func (c *Controller) GetBlocks(ctx *fiber.Ctx) error {
	q := new(BlockQuery)
	err := ctx.QueryParser(q)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error":      err.Error(),
		})
	}

	b, err := c.s.GetBlocks(q)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"successful": false,
			"error":      err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"successful": true,
		"count":      b.Count,
		"blocks":     b,
	})
}

func (c *Controller) GetAddress(ctx *fiber.Ctx) error {
	pub := ctx.Query("pub")
	if len(pub) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error":      "public key not provided",
		})
	}

	addr, err := c.s.GetAddress(pub)
	if err != nil {
		if err == ErrInvalidPubKey {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"successful": false,
				"error":      "invalid public key",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"successful": false,
			"error":      err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"successful": true,
		"public_key": pub,
		"address":    addr,
	})
}

func (c *Controller) GetMyTXs(ctx *fiber.Ctx) error {
	addr := ctx.Query("addr")
	if len(addr) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error":      "address not provided",
		})
	}

	txs, err := c.s.GetMyTXs(addr)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"successful": false,
			"error":      err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"successful": true,
		"txs":        txs,
	})
}
