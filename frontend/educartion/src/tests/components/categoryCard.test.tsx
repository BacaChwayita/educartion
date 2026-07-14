import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import CategoryCard from '@/components/CategoryCard'

const mockCategory = {
  category_id: 1,
  name: 'Electronics',
  description: 'Electronic products',
  is_active: true,
}

describe('CategoryCard', () => {
  test('renders category name', () => {
    render(<CategoryCard category={mockCategory} />)
    expect(screen.getByText('Electronics')).toBeInTheDocument()
  })

  test('renders category description', () => {
    render(<CategoryCard category={mockCategory} />)
    expect(screen.getByText('Electronic products')).toBeInTheDocument()
  })

  test('does not render inactive category', () => {
    render(<CategoryCard category={{ ...mockCategory, is_active: false }} />)
    expect(screen.queryByText('Electronics')).not.toBeInTheDocument()
  })
})