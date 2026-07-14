import type { Metadata } from "next";
import { AuthForm } from "@/components/auth-form";

export const metadata: Metadata = {
  title: "Register | educartion",
  description: "Create an educartion account.",
};

export default function RegisterPage() {
  return <AuthForm mode="register" />;
}