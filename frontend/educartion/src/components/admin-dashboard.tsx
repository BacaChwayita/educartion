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
    return <div>Loading...</div>;
  }

  return (
    <main style={{ width: "100%", padding: 24, display: "grid", justifyItems: "center" }}>
      <div style={{ width: "100%", maxWidth: 1200, display: "grid", gap: 24 }}>
        <header>
          <h1 style={{ margin: 0, fontSize: 32 }}>Admin Dashboard</h1>
          <p style={{ margin: "8px 0 0", color: "#555" }}>{message}</p>
        </header>

        <section style={{ padding: 24, border: "1px solid #e5e7eb", borderRadius: 16, background: "white" }}>
          <h2 style={{ margin: 0, marginBottom: 16 }}>Dashboard Overview</h2>
          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(250px, 1fr))", gap: 16 }}>
            <div style={{ padding: 16, border: "1px solid #e5e7eb", borderRadius: 12, background: "#f9fafb" }}>
              <h3 style={{ margin: 0, fontSize: 14, color: "#6b7280" }}>Total Users</h3>
              <p style={{ margin: "8px 0 0", fontSize: 28, fontWeight: 700 }}>—</p>
            </div>
            <div style={{ padding: 16, border: "1px solid #e5e7eb", borderRadius: 12, background: "#f9fafb" }}>
              <h3 style={{ margin: 0, fontSize: 14, color: "#6b7280" }}>Total Orders</h3>
              <p style={{ margin: "8px 0 0", fontSize: 28, fontWeight: 700 }}>—</p>
            </div>
            <div style={{ padding: 16, border: "1px solid #e5e7eb", borderRadius: 12, background: "#f9fafb" }}>
              <h3 style={{ margin: 0, fontSize: 14, color: "#6b7280" }}>Total Revenue</h3>
              <p style={{ margin: "8px 0 0", fontSize: 28, fontWeight: 700 }}>—</p>
            </div>
          </div>
        </section>

        <section style={{ padding: 24, border: "1px solid #e5e7eb", borderRadius: 16, background: "white" }}>
          <h2 style={{ margin: 0, marginBottom: 16 }}>Recent Activity</h2>
          <p style={{ margin: 0, color: "#6b7280" }}>No recent activity to display.</p>
        </section>
      </div>
    </main>
  );
}
