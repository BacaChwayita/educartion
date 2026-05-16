package cart

import (
	"errors"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func GetCartByToken(token string) (model.Cart, error) {
	cart, err := db.GetCartByJTI(token)
	if err == nil {
		return cart, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return model.Cart{}, err
	}

	account, err := db.GetAccountByJTI(token)
	if err != nil {
		return model.Cart{}, err
	}

	newCart := model.Cart{
		Account_id:  account.Account_id,
		Session_key: token,
		Total_Price: 0,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	return db.InsertCart(newCart)
}

func GetCartItems(token string) ([]model.Cart_item, error) {
	cart, err := GetCartByToken(token)
	if err != nil {
		return nil, err
	}

	return db.GetCartItemsByCart(cart.Cart_id)
}

func AddCartItem(token string, req model.AddCartItemRequest) (model.Cart_item, error) {
	if req.Quantity <= 0 {
		return model.Cart_item{}, ErrInvalidQuantity
	}

	cart, err := GetCartByToken(token)
	if err != nil {
		return model.Cart_item{}, err
	}

	product, err := db.GetProductById(req.ProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Cart_item{}, ErrProductNotFound
		}
		return model.Cart_item{}, err
	}

	cartItem, err := db.GetCartItem(cart.Cart_id, req.ProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			newItem := model.Cart_item{
				Cart_id:    cart.Cart_id,
				Product_id: req.ProductID,
				Quantity:   req.Quantity,
				Unit_price: product.Price,
			}
			cartItem, err = db.InsertCartItem(newItem)
			if err != nil {
				return model.Cart_item{}, err
			}
		} else {
			return model.Cart_item{}, err
		}
	} else {
		cartItem, err = db.UpdateCartItem(cart.Cart_id, req.ProductID, cartItem.Quantity+req.Quantity)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.Cart_item{}, ErrCartItemNotFound
			}
			return model.Cart_item{}, err
		}
	}

	if err := refreshCartTotal(&cart); err != nil {
		return model.Cart_item{}, err
	}

	return cartItem, nil
}

func UpdateCartItem(token string, productID int, quantity int) (model.Cart_item, error) {
	if quantity <= 0 {
		return model.Cart_item{}, ErrInvalidQuantity
	}

	cart, err := GetCartByToken(token)
	if err != nil {
		return model.Cart_item{}, err
	}

	updatedItem, err := db.UpdateCartItem(cart.Cart_id, productID, quantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Cart_item{}, ErrCartItemNotFound
		}
		return model.Cart_item{}, err
	}

	if err := refreshCartTotal(&cart); err != nil {
		return model.Cart_item{}, err
	}

	return updatedItem, nil
}

func DeleteCartItem(token string, productID int) error {
	cart, err := GetCartByToken(token)
	if err != nil {
		return err
	}

	rows, err := db.DeleteCartItem(cart.Cart_id, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCartItemNotFound
		}
		return err
	}

	if rows == 0 {
		return ErrCartItemNotFound
	}

	return refreshCartTotal(&cart)
}

func refreshCartTotal(cart *model.Cart) error {
	items, err := db.GetCartItemsByCart(cart.Cart_id)
	if err != nil {
		return err
	}

	total := 0
	for _, item := range items {
		total += item.Quantity * item.Unit_price
	}

	cart.Total_Price = total
	cart.Updated_at = time.Now()

	_, err = db.UpdateCart(*cart)
	return err
}
