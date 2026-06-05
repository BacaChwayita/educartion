type CartItemProps = {
  cart_id: number;
  product_id: number;
  quantity: number;
  unit_price: number;
};

type Props = {
  item: CartItemProps;
  onRemove: (productId: number) => void;
  onUpdate: (productId: number, quantity: number) => void;
};

export default function CartItem({ item, onRemove, onUpdate }: Props) {
  return (
    <article>
      <div>
        Quantity: <span>{item.quantity}</span>
      </div>
      <div>
        Price: <span>{item.unit_price}</span>
      </div>
      <button type="button" onClick={() => onUpdate(item.product_id, item.quantity + 1)}>
        Increase quantity
      </button>
      <button type="button" onClick={() => onUpdate(item.product_id, item.quantity - 1)}>
        Decrease quantity
      </button>
      <button type="button" onClick={() => onRemove(item.product_id)}>
        Remove from cart
      </button>
    </article>
  );
}
