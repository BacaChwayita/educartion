"use client";

import React, { useEffect, useState } from "react";
import * as adminController from "@/controllers/adminController";

export function AdminDashboard() {
  const [loaded, setLoaded] = useState(false);
  const [message, setMessage] = useState("");

  useEffect(() => {
    const loadDashboard = async () => {
      const result = await adminController.handleLoadAdminDashboard();
      if (result.ok) {
        setMessage(result.data.message);
      }
      setLoaded(true);
    };
    loadDashboard();
  }, []);

  if (!loaded) {
    return (
      <section className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,rgba(251,191,36,0.2),transparent_30%),radial-gradient(circle_at_bottom_right,rgba(14,165,233,0.18),transparent_28%),linear-gradient(180deg,#060816_0%,#0b1020_100%)] px-6 py-10 text-slate-100 sm:px-8 lg:px-10">
        <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-size-[28px_28px] opacity-20" />
        <div className="relative mx-auto flex min-h-[calc(100vh-5rem)] w-full items-center justify-center">
          <p className="text-slate-300">Loading dashboard...</p>
        </div>
      </section>
    );
  }

  return (
    <section className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,rgba(251,191,36,0.2),transparent_30%),radial-gradient(circle_at_bottom_right,rgba(14,165,233,0.18),transparent_28%),linear-gradient(180deg,#060816_0%,#0b1020_100%)] px-6 py-10 text-slate-100 sm:px-8 lg:px-10">
      <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-size-[28px_28px] opacity-20" />
      <div className="relative mx-auto w-full max-w-6xl">
        <div className="w-full rounded-4xl border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30 backdrop-blur-xl lg:p-8">
          <div className="rounded-3xl border border-white/10 bg-slate-950/80 p-6 sm:p-8">
            <div className="mb-8 space-y-2">
              <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-200">
                Administration
              </p>
              <h2 className="text-3xl font-semibold tracking-tight text-white">Admin Dashboard</h2>
              <p className="text-sm leading-6 text-slate-300">{message}</p>
            </div>

            <div className="space-y-6">
              <div className="rounded-2xl border border-white/10 bg-white/5 p-6">
                <h3 className="text-sm font-semibold text-white mb-4 uppercase tracking-[0.2em] text-amber-200">Dashboard Overview</h3>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="rounded-xl border border-white/10 bg-slate-900/30 p-4">
                    <h4 className="text-xs font-semibold text-slate-400 mb-2 uppercase">Total Users</h4>
                    <p className="text-2xl font-bold text-amber-400">—</p>
                  </div>
                  <div className="rounded-xl border border-white/10 bg-slate-900/30 p-4">
                    <h4 className="text-xs font-semibold text-slate-400 mb-2 uppercase">Total Orders</h4>
                    <p className="text-2xl font-bold text-amber-400">—</p>
                  </div>
                  <div className="rounded-xl border border-white/10 bg-slate-900/30 p-4">
                    <h4 className="text-xs font-semibold text-slate-400 mb-2 uppercase">Total Revenue</h4>
                    <p className="text-2xl font-bold text-amber-400">—</p>
                  </div>
                </div>
              </div>

              <div className="rounded-2xl border border-white/10 bg-white/5 p-6">
                <h3 className="text-sm font-semibold text-white mb-4 uppercase tracking-[0.2em] text-amber-200">Recent Activity</h3>
                <p className="text-slate-400">No recent activity to display.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
