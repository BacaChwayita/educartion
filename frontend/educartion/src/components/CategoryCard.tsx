type Category = {
  category_id: number;
  name: string;
  description: string;
  is_active: boolean;
};

export default function CategoryCard({ category }: { category: Category }) {
  if (!category.is_active) {
    return null;
  }

  return (
    <article>
      <h2>{category.name}</h2>
      <p>{category.description}</p>
    </article>
  );
}
