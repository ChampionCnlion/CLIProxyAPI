package usagemonitor

const PanelHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>CPA Request Monitor</title>
  <style>
    :root {
      --bg: #f5efe4;
      --ink: #202018;
      --muted: #6f6a5f;
      --line: #d8ccb8;
      --card: rgba(255, 252, 245, 0.92);
      --accent: #2f6b55;
      --accent-2: #b95f37;
      --bad: #a33b31;
      --good: #247350;
      --shadow: 0 18px 60px rgba(59, 44, 23, 0.14);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      color: var(--ink);
      font-family: ui-serif, Georgia, Cambria, "Times New Roman", serif;
      background:
        radial-gradient(circle at top left, rgba(47, 107, 85, 0.24), transparent 34rem),
        radial-gradient(circle at bottom right, rgba(185, 95, 55, 0.18), transparent 32rem),
        linear-gradient(135deg, #f8f1e6, var(--bg));
    }
    main { width: min(1180px, calc(100% - 32px)); margin: 0 auto; padding: 38px 0 56px; }
    header { display: flex; justify-content: space-between; gap: 22px; align-items: flex-start; margin-bottom: 22px; }
    h1 { margin: 0; font-size: clamp(2.2rem, 5vw, 4.6rem); letter-spacing: -0.06em; line-height: 0.95; }
    .sub { margin-top: 12px; color: var(--muted); max-width: 760px; font-size: 1rem; line-height: 1.6; }
    .auth {
      min-width: min(100%, 420px);
      padding: 16px;
      border: 1px solid var(--line);
      background: var(--card);
      box-shadow: var(--shadow);
      border-radius: 22px;
    }
    .auth label { display: block; color: var(--muted); font-size: 0.82rem; margin-bottom: 8px; text-transform: uppercase; letter-spacing: 0.08em; }
    .row { display: flex; gap: 10px; }
    input, button {
      border: 1px solid var(--line);
      border-radius: 14px;
      padding: 11px 13px;
      font: inherit;
      background: #fffaf0;
      color: var(--ink);
    }
    input { min-width: 0; flex: 1; }
    button { cursor: pointer; background: var(--accent); color: white; border-color: var(--accent); }
    button.secondary { background: #fffaf0; color: var(--ink); }
    button:disabled { opacity: 0.55; cursor: wait; }
    .grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 14px; margin: 18px 0; }
    .card {
      border: 1px solid var(--line);
      background: var(--card);
      border-radius: 22px;
      padding: 18px;
      box-shadow: var(--shadow);
    }
    .metric span { display: block; color: var(--muted); font-size: 0.82rem; text-transform: uppercase; letter-spacing: 0.08em; }
    .metric strong { display: block; margin-top: 8px; font-size: clamp(1.8rem, 4vw, 3rem); letter-spacing: -0.04em; }
    .toolbar { display: flex; justify-content: space-between; align-items: center; gap: 14px; margin: 22px 0 12px; }
    .status { color: var(--muted); font-size: 0.94rem; }
    .tableWrap { overflow: auto; border: 1px solid var(--line); border-radius: 22px; background: var(--card); box-shadow: var(--shadow); }
    table { border-collapse: collapse; width: 100%; min-width: 920px; }
    th, td { text-align: left; padding: 12px 14px; border-bottom: 1px solid rgba(216, 204, 184, 0.72); vertical-align: top; }
    th { color: var(--muted); font-size: 0.78rem; text-transform: uppercase; letter-spacing: 0.08em; background: rgba(255, 250, 240, 0.74); position: sticky; top: 0; }
    tr:last-child td { border-bottom: 0; }
    .pill { display: inline-flex; padding: 4px 9px; border-radius: 999px; font-size: 0.78rem; border: 1px solid var(--line); }
    .ok { color: var(--good); border-color: rgba(36, 115, 80, 0.32); background: rgba(36, 115, 80, 0.08); }
    .fail { color: var(--bad); border-color: rgba(163, 59, 49, 0.32); background: rgba(163, 59, 49, 0.08); }
    .muted { color: var(--muted); }
    .error { color: var(--bad); white-space: pre-wrap; }
    @media (max-width: 860px) {
      header { display: block; }
      .auth { margin-top: 18px; }
      .grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
      .row { flex-wrap: wrap; }
      button { flex: 1; }
    }
  </style>
</head>
<body>
  <main>
    <header>
      <section>
        <h1>Request Monitor</h1>
        <div class="sub">Built into CLIProxyAPI. Data is persisted locally in SQLite and read from the same management API, so a separate CPA-Manager container is not required.</div>
      </section>
      <section class="auth">
        <label for="key">Management Key</label>
        <div class="row">
          <input id="key" type="password" autocomplete="current-password" placeholder="Bearer token for /v0/management">
          <button id="saveKey">Save</button>
          <button id="clearKey" class="secondary">Clear</button>
        </div>
      </section>
    </header>

    <section class="grid">
      <div class="card metric"><span>Total requests</span><strong id="totalRequests">-</strong></div>
      <div class="card metric"><span>Success</span><strong id="successCount">-</strong></div>
      <div class="card metric"><span>Failure</span><strong id="failureCount">-</strong></div>
      <div class="card metric"><span>Total tokens</span><strong id="totalTokens">-</strong></div>
    </section>

    <section class="card">
      <div id="serviceStatus" class="status">Status not loaded.</div>
      <div id="error" class="error"></div>
    </section>

    <div class="toolbar">
      <div class="status" id="rowCount">No rows loaded.</div>
      <div class="row">
        <button id="refresh">Refresh</button>
        <button id="exportBtn" class="secondary">Export JSONL</button>
      </div>
    </div>

    <section class="tableWrap">
      <table>
        <thead>
          <tr>
            <th>Time</th>
            <th>Endpoint</th>
            <th>Model</th>
            <th>Source</th>
            <th>Auth</th>
            <th>Status</th>
            <th>Tokens</th>
            <th>Latency</th>
          </tr>
        </thead>
        <tbody id="rows"></tbody>
      </table>
    </section>
  </main>
  <script>
    const keyInput = document.getElementById('key');
    const saveKey = document.getElementById('saveKey');
    const clearKey = document.getElementById('clearKey');
    const refreshBtn = document.getElementById('refresh');
    const exportBtn = document.getElementById('exportBtn');
    const errorBox = document.getElementById('error');
    const rowsEl = document.getElementById('rows');
    const rowCount = document.getElementById('rowCount');
    const serviceStatus = document.getElementById('serviceStatus');
    const stored = sessionStorage.getItem('cpa_usage_management_key') || '';
    keyInput.value = stored;

    saveKey.onclick = () => {
      sessionStorage.setItem('cpa_usage_management_key', keyInput.value.trim());
      load();
    };
    clearKey.onclick = () => {
      keyInput.value = '';
      sessionStorage.removeItem('cpa_usage_management_key');
    };
    refreshBtn.onclick = () => load();
    exportBtn.onclick = () => {
      const headers = authHeaders();
      fetch('/v0/management/usage/export', { headers }).then(async res => {
        if (!res.ok) throw new Error(await res.text());
        return res.blob();
      }).then(blob => {
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'usage-events.jsonl';
        a.click();
        URL.revokeObjectURL(url);
      }).catch(showError);
    };

    function authHeaders() {
      const token = keyInput.value.trim() || sessionStorage.getItem('cpa_usage_management_key') || '';
      return token ? { Authorization: 'Bearer ' + token } : {};
    }
    async function fetchJSON(path) {
      const res = await fetch(path, { headers: authHeaders() });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(path + ' failed: ' + res.status + ' ' + text);
      }
      return res.json();
    }
    function fmt(n) {
      if (n === undefined || n === null) return '-';
      return Number(n).toLocaleString();
    }
    function flattenUsage(payload) {
      const out = [];
      const apis = payload.apis || {};
      for (const endpoint of Object.keys(apis)) {
        const models = (apis[endpoint] && apis[endpoint].models) || {};
        for (const model of Object.keys(models)) {
          const details = models[model].details || [];
          for (const detail of details) out.push({ endpoint, model, detail });
        }
      }
      out.sort((a, b) => new Date(b.detail.timestamp) - new Date(a.detail.timestamp));
      return out;
    }
    function renderRows(items) {
      rowsEl.innerHTML = '';
      for (const item of items.slice(0, 500)) {
        const d = item.detail || {};
        const t = d.tokens || {};
        const tr = document.createElement('tr');
        tr.innerHTML =
          '<td>' + esc(new Date(d.timestamp).toLocaleString()) + '</td>' +
          '<td>' + esc(item.endpoint || '-') + '</td>' +
          '<td>' + esc(item.model || '-') + '</td>' +
          '<td>' + esc(d.source || '-') + '</td>' +
          '<td>' + esc(d.auth_index || '-') + '</td>' +
          '<td><span class="pill ' + (d.failed ? 'fail' : 'ok') + '">' + (d.failed ? 'failed' : 'ok') + '</span></td>' +
          '<td>' + fmt(t.total_tokens) + ' <span class="muted">in ' + fmt(t.input_tokens) + ' / out ' + fmt(t.output_tokens) + '</span></td>' +
          '<td>' + (d.latency_ms === undefined ? '-' : fmt(d.latency_ms) + ' ms') + '</td>';
        rowsEl.appendChild(tr);
      }
      rowCount.textContent = items.length.toLocaleString() + ' usage detail rows loaded. Showing latest 500.';
    }
    function esc(value) {
      return String(value).replace(/[&<>"']/g, ch => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[ch]));
    }
    function showError(err) {
      errorBox.textContent = err && err.message ? err.message : String(err);
    }
    async function load() {
      errorBox.textContent = '';
      refreshBtn.disabled = true;
      try {
        const [status, usage] = await Promise.all([fetchJSON('/status'), fetchJSON('/v0/management/usage')]);
        document.getElementById('totalRequests').textContent = fmt(usage.total_requests);
        document.getElementById('successCount').textContent = fmt(usage.success_count);
        document.getElementById('failureCount').textContent = fmt(usage.failure_count);
        document.getElementById('totalTokens').textContent = fmt(usage.total_tokens);
        const c = status.collector || {};
        serviceStatus.textContent = 'collector=' + (c.collector || '-') + ', transport=' + (c.transport || '-') + ', db=' + (status.dbPath || '-') + ', events=' + fmt(status.events) + ', deadLetters=' + fmt(status.deadLetters) + (c.lastError ? ', lastError=' + c.lastError : '');
        renderRows(flattenUsage(usage));
      } catch (err) {
        showError(err);
      } finally {
        refreshBtn.disabled = false;
      }
    }
    load();
  </script>
</body>
</html>`
