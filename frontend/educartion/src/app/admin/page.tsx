import type { Metadata } from "next";
import { AdminDashboard } from "@/components/admin-dashboard";

export const metadata: Metadata = {
  title: "Admin | educartion",
  description: "Admin dashboard for educartion.",
};

export default function AdminPage() {
  return <AdminDashboard />;
}
