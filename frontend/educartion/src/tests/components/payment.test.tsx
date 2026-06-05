import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import PaymentPage from '@/app/payment/page'

const pushMock = jest.fn()

jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: pushMock,
  }),
}))

describe('PaymentPage', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    window.localStorage.clear()
  })

  it('renders empty cart state when there are no items', async () => {
    window.localStorage.setItem('cart', JSON.stringify([]))

    render(<PaymentPage />)

    expect(await screen.findByText('Your cart is empty.')).toBeInTheDocument()
    expect(screen.getByText('Subtotal')).toBeInTheDocument()
  })

  it('displays cart items and totals from localStorage', async () => {
    window.localStorage.setItem(
      'cart',
      JSON.stringify([
        { id: '1', title: 'Test product', price: 10, qty: 2 },
      ])
    )

    render(<PaymentPage />)

    expect(await screen.findByText('Test product')).toBeInTheDocument()
    expect(screen.getAllByText('R20.00').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('R5.00')).toBeInTheDocument()
    expect(screen.getByText('R25.00')).toBeInTheDocument()
  })

  it('shows an alert when checkout details are incomplete', async () => {
    const alertMock = jest.spyOn(window, 'alert').mockImplementation(() => {})
    window.localStorage.setItem(
      'cart',
      JSON.stringify([{ id: '1', title: 'Test product', price: 10, qty: 1 }])
    )

    const user = userEvent.setup()
    render(<PaymentPage />)

    await screen.findByText('Test product')
    await user.click(screen.getByRole('button', { name: /confirm order/i }))

    expect(alertMock).toHaveBeenCalledWith('Please complete your delivery and payment details.')
    expect(pushMock).not.toHaveBeenCalled()

    alertMock.mockRestore()
  })

  it('completes checkout and navigates to orders when input is valid', async () => {
    const alertMock = jest.spyOn(window, 'alert').mockImplementation(() => {})
    window.localStorage.setItem(
      'cart',
      JSON.stringify([{ id: '1', title: 'Test product', price: 10, qty: 1 }])
    )

    const user = userEvent.setup()
    render(<PaymentPage />)

    await screen.findByText('Test product')

    await user.type(screen.getByPlaceholderText('Full name'), 'John Doe')
    await user.type(screen.getByPlaceholderText('Street address'), '123 Main St')
    await user.type(screen.getByPlaceholderText('City'), 'Cape Town')
    await user.type(screen.getByPlaceholderText('Postal code'), '8000')
    await user.type(screen.getByPlaceholderText('Name on card'), 'John Doe')
    await user.type(screen.getByPlaceholderText('Card number'), '123456789012')
    await user.type(screen.getByPlaceholderText('MM/YY'), '12/25')
    await user.type(screen.getByPlaceholderText('CVC'), '123')

    await user.click(screen.getByRole('button', { name: /confirm order/i }))

    await waitFor(() => {
      expect(pushMock).toHaveBeenCalledWith('/orders')
    })
    expect(window.localStorage.getItem('cart')).toBeNull()

    alertMock.mockRestore()
  })
})
