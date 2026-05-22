"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import {
  type LoginClientRequest,
  type LoginResponse,
  type RegisterClientRequest,
} from "@/lib/auth-contract";
import { readBackendError } from "@/lib/auth-errors";

type AuthMode = "login" | "register";

interface AuthFormProps {
  mode: AuthMode;
}

interface ValidationState {
  error: string | null;
  status: string | null;
}

const initialValidationState: ValidationState = {
  error: null,
  status: null,
};

export function AuthForm({ mode }: AuthFormProps) {
  const router = useRouter();
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(false);
  const [formState, setFormState] = useState<ValidationState>(initialValidationState);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const isLogin = mode === "login";
  const title = isLogin ? "Sign in" : "Create your account";
  const description = isLogin
    ? "Use your email and password to return to your account."
    : "Register once and use the same account for the whole platform.";
  const panelHeadline = isLogin
    ? "Welcome back to educartion"
    : "Start your educartion workspace";
  const panelCopy = isLogin
    ? "Pick up where you left off and keep your session alive when you choose remember me."
    : "Set up your profile in a single pass. The backend expects the same snake_case payload fields it already uses.";

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFormState(initialValidationState);

    if (isLogin) {
      const payload: LoginClientRequest = {
        email: email.trim(),
        password,
        rememberMe,
      };

      if (!payload.email || !payload.password) {
        setFormState({
          error: "Email and password are required.",
          status: null,
        });
        return;
      }

      setIsSubmitting(true);

      try {
        const response = await fetch("/api/auth/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(payload),
        });

        if (!response.ok) {
          setFormState({
            error: await readBackendError(response),
            status: null,
          });
          return;
        }

        const data = (await response.json()) as LoginResponse;
        setFormState({
          error: null,
          status: `Signed in as ${data.full_name} (${data.role}).`,
        });
      } catch {
        setFormState({
          error: "Unable to reach the login endpoint right now.",
          status: null,
        });
      } finally {
        setIsSubmitting(false);
      }

      return;
    }

    const trimmedFullName = fullName.trim();
    const trimmedEmail = email.trim();

    const payload: RegisterClientRequest = {
      fullName: trimmedFullName,
      email: trimmedEmail,
      password,
    };

    if (!payload.fullName || !payload.email || !payload.password) {
      setFormState({
        error: "Full name, email, and password are required.",
        status: null,
      });
      return;
    }

    if (password !== confirmPassword) {
      setFormState({
        error: "Passwords do not match.",
        status: null,
      });
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await fetch("/api/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        setFormState({
          error: await readBackendError(response),
          status: null,
        });
        return;
      }

      setFullName("");
      setEmail("");
      setPassword("");
      setConfirmPassword("");
      setRememberMe(false);
      setFormState({
        error: null,
        status: "Account created. You can sign in now.",
      });
      router.push("/login");
    } catch {
      setFormState({
        error: "Unable to reach the register endpoint right now.",
        status: null,
      });
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <section className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,_rgba(251,191,36,0.2),_transparent_30%),radial-gradient(circle_at_bottom_right,_rgba(14,165,233,0.18),_transparent_28%),linear-gradient(180deg,#060816_0%,#0b1020_100%)] px-6 py-10 text-slate-100 sm:px-8 lg:px-10">
      <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:28px_28px] opacity-20" />
      <div className="relative mx-auto flex min-h-[calc(100vh-5rem)] w-full max-w-6xl items-center">
        <div className="grid w-full gap-8 rounded-[2rem] border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30 backdrop-blur-xl lg:grid-cols-[1.1fr_0.9fr] lg:p-8">
          <aside className="flex flex-col justify-between rounded-[1.5rem] border border-white/10 bg-slate-950/55 p-6">
            <div className="space-y-6">
              <div className="inline-flex items-center gap-3 rounded-full border border-amber-400/30 bg-amber-400/10 px-4 py-2 text-sm font-medium text-amber-200">
                <span className="h-2.5 w-2.5 rounded-full bg-amber-300" />
                educartion auth
              </div>
              <div className="space-y-4">
                <p className="text-sm uppercase tracking-[0.35em] text-slate-400">
                  Secure access
                </p>
                <h1 className="max-w-xl text-4xl font-semibold tracking-tight text-white sm:text-5xl">
                  {panelHeadline}
                </h1>
                <p className="max-w-xl text-base leading-7 text-slate-300">
                  {panelCopy}
                </p>
              </div>
            </div>

            <div className="mt-10 grid gap-3 sm:grid-cols-3">
              <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                <p className="text-xs uppercase tracking-[0.25em] text-slate-400">
                  Endpoint
                </p>
                <p className="mt-2 text-sm text-slate-100">/api/auth/{isLogin ? "login" : "register"}</p>
              </div>
              <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                <p className="text-xs uppercase tracking-[0.25em] text-slate-400">
                  Payload
                </p>
                <p className="mt-2 text-sm text-slate-100">snake_case server contract</p>
              </div>
              <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                <p className="text-xs uppercase tracking-[0.25em] text-slate-400">
                  Session
                </p>
                <p className="mt-2 text-sm text-slate-100">remember me cookie support</p>
              </div>
            </div>
          </aside>

          <div className="rounded-[1.5rem] border border-white/10 bg-slate-950/80 p-6 sm:p-8">
            <div className="mb-8 space-y-3">
              <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-200">
                {isLogin ? "Login" : "Register"}
              </p>
              <h2 className="text-3xl font-semibold tracking-tight text-white">{title}</h2>
              <p className="text-sm leading-6 text-slate-300">{description}</p>
            </div>

            <form className="space-y-5" onSubmit={handleSubmit}>
              {!isLogin ? (
                <label className="block space-y-2">
                  <span className="text-sm font-medium text-slate-200">Full name</span>
                  <input
                    name="fullName"
                    autoComplete="name"
                    value={fullName}
                    onChange={(event) => setFullName(event.target.value)}
                    placeholder="Jane Doe"
                    className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                  />
                </label>
              ) : null}

              <label className="block space-y-2">
                <span className="text-sm font-medium text-slate-200">Email</span>
                <input
                  name="email"
                  type="email"
                  autoComplete="email"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  placeholder="you@example.com"
                  className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                />
              </label>

              <label className="block space-y-2">
                <span className="text-sm font-medium text-slate-200">Password</span>
                <input
                  name="password"
                  type="password"
                  autoComplete={isLogin ? "current-password" : "new-password"}
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  placeholder="••••••••"
                  className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                />
              </label>

              {!isLogin ? (
                <label className="block space-y-2">
                  <span className="text-sm font-medium text-slate-200">Confirm password</span>
                  <input
                    name="confirmPassword"
                    type="password"
                    autoComplete="new-password"
                    value={confirmPassword}
                    onChange={(event) => setConfirmPassword(event.target.value)}
                    placeholder="••••••••"
                    className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                  />
                </label>
              ) : (
                <label className="flex items-center gap-3 rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-slate-200">
                  <input
                    name="rememberMe"
                    type="checkbox"
                    checked={rememberMe}
                    onChange={(event) => setRememberMe(event.target.checked)}
                    className="h-4 w-4 rounded border-white/20 bg-white/10 text-amber-400 focus:ring-amber-300/20"
                  />
                  Remember me on this device
                </label>
              )}

              {formState.error ? (
                <div className="rounded-2xl border border-rose-400/30 bg-rose-400/10 px-4 py-3 text-sm text-rose-100">
                  {formState.error}
                </div>
              ) : null}

              {formState.status ? (
                <div className="rounded-2xl border border-emerald-400/30 bg-emerald-400/10 px-4 py-3 text-sm text-emerald-50">
                  {formState.status}
                </div>
              ) : null}

              <button
                type="submit"
                disabled={isSubmitting}
                className="flex w-full items-center justify-center rounded-2xl bg-amber-400 px-5 py-3.5 text-sm font-semibold text-slate-950 transition hover:bg-amber-300 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {isSubmitting ? "Please wait..." : isLogin ? "Login" : "Register"}
              </button>
            </form>

            <div className="mt-6 rounded-2xl border border-white/10 bg-white/5 px-4 py-4 text-sm text-slate-300">
              <div className="flex flex-wrap items-center gap-3">
                <span>{isLogin ? "Don't have an account?" : "Already have an account?"}</span>
                <Link
                  href={isLogin ? "/register" : "/login"}
                  className="inline-flex items-center rounded-full border border-white/15 px-4 py-2 font-medium text-white transition hover:bg-white/10"
                >
                  {isLogin ? "Register" : "Login"}
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
