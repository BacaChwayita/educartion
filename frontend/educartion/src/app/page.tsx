import Link from "next/link";

export default function Home() {
  return (
    <section className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,rgba(251,191,36,0.2),transparent_30%),radial-gradient(circle_at_bottom_right,rgba(14,165,233,0.18),transparent_28%),linear-gradient(180deg,#060816_0%,#0b1020_100%)] px-6 py-10 text-slate-100 sm:px-8 lg:px-10">
      <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-size-[28px_28px] opacity-20" />
      <div className="relative mx-auto flex min-h-[calc(100vh-5rem)] w-full max-w-4xl items-center justify-center">
        <div className="w-full rounded-4xl border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30 backdrop-blur-xl lg:p-8">
          <div className="rounded-3xl border border-white/10 bg-slate-950/80 p-6 sm:p-8">
            <div className="mb-8 space-y-3">
              <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-200">
                Welcome
              </p>
              <h2 className="text-4xl font-semibold tracking-tight text-white">
                educartion
              </h2>
              <p className="text-base leading-6 text-slate-300">
                A comprehensive platform for education management and learning.
              </p>
            </div>

            <div className="mb-8 space-y-4">
              <p className="text-sm text-slate-400">
                Get started by logging in to your account or creating a new one.
              </p>
              <div className="flex flex-col gap-3 sm:flex-row">
                <Link
                  href="/login"
                  className="flex items-center justify-center rounded-2xl bg-amber-400 px-6 py-3 text-sm font-semibold text-slate-950 transition hover:bg-amber-300"
                >
                  Login
                </Link>
                <Link
                  href="/register"
                  className="flex items-center justify-center rounded-2xl border border-white/15 px-6 py-3 text-sm font-semibold text-white transition hover:bg-white/10"
                >
                  Register
                </Link>
              </div>
            </div>

            <div className="rounded-2xl border border-white/10 bg-white/5 px-4 py-4">
              <p className="text-xs font-medium uppercase tracking-[0.2em] text-amber-200 mb-3">
                Features
              </p>
              <ul className="space-y-2 text-sm text-slate-300">
                <li className="flex items-start gap-3">
                  <span className="mt-1 h-1.5 w-1.5 rounded-full bg-amber-400 flex-shrink-0" />
                  <span>Secure user authentication and account management</span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="mt-1 h-1.5 w-1.5 rounded-full bg-amber-400 flex-shrink-0" />
                  <span>Comprehensive course and education catalog</span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="mt-1 h-1.5 w-1.5 rounded-full bg-amber-400 flex-shrink-0" />
                  <span>Order management and payment processing</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
