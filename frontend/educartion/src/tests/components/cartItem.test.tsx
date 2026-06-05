import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import CartItem from '@/components/CartItem'

const mockItem = {
  cart_id: 1,
  product_id: 1,
  quantity: 2,
  unit_price: 1999,
}

describe('CartItem', () => {
  test('renders quantity correctly', () => {
    render(<CartItem item={mockItem} onRemove={jest.fn()} onUpdate={jest.fn()} />)
    expect(screen.getByText('2')).toBeInTheDocument()
  })

  test('renders unit price correctly', () => {
    render(<CartItem item={mockItem} onRemove={jest.fn()} onUpdate={jest.fn()} />)
    expect(screen.getByText('1999')).toBeInTheDocument()
  })

  test('calls onRemove when delete button is clicked', async () => {
    const onRemove = jest.fn()
    render(<CartItem item={mockItem} onRemove={onRemove} onUpdate={jest.fn()} />)

    await userEvent.click(screen.getByRole('button', { name: /remove/i }))
    expect(onRemove).toHaveBeenCalledTimes(1)
    expect(onRemove).toHaveBeenCalledWith(mockItem.product_id)
  })

  test('calls onUpdate with increased quantity when increase button is clicked', async () => {
    const onUpdate = jest.fn()
    render(<CartItem item={mockItem} onRemove={jest.fn()} onUpdate={onUpdate} />)

    await userEvent.click(screen.getByRole('button', { name: /increase/i }))
    expect(onUpdate).toHaveBeenCalledWith(mockItem.product_id, 3)
  })

  test('calls onUpdate with decreased quantity when decrease button is clicked', async () => {
    const onUpdate = jest.fn()
    render(<CartItem item={mockItem} onRemove={jest.fn()} onUpdate={onUpdate} />)

    await userEvent.click(screen.getByRole('button', { name: /decrease/i }))
    expect(onUpdate).toHaveBeenCalledWith(mockItem.product_id, 1)
  })
})