package usagemonitor

const PanelHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>CPA Request Monitor</title>
  <style>
    :root {
      --bg: #eef1e8;
      --ink: #17201c;
      --muted: #66736c;
      --line: rgba(36, 52, 45, 0.16);
      --card: rgba(255, 253, 246, 0.9);
      --card-strong: #fffdf6;
      --accent: #255b49;
      --accent-2: #c46f38;
      --good: #19764a;
      --warn: #b66a13;
      --bad: #b13d35;
      --shadow: 0 24px 80px rgba(33, 47, 39, 0.14);
      --mono: "SFMono-Regular", Consolas, "Liberation Mono", monospace;
    }
    * { box-sizing: border-box; }
    html { color-scheme: light; }
    body {
      margin: 0;
      min-height: 100vh;
      color: var(--ink);
      font-family: Athelas, "Iowan Old Style", Georgia, Cambria, "Times New Roman", serif;
      background:
        radial-gradient(circle at 10% 8%, rgba(37, 91, 73, 0.26), transparent 24rem),
        radial-gradient(circle at 86% 12%, rgba(196, 111, 56, 0.22), transparent 26rem),
        linear-gradient(135deg, #f7f0df 0%, #edf2e9 42%, #e8efe8 100%);
    }
    body::before {
      content: "";
      position: fixed;
      inset: 0;
      pointer-events: none;
      opacity: 0.28;
      background-image:
        linear-gradient(rgba(23, 32, 28, 0.05) 1px, transparent 1px),
        linear-gradient(90deg, rgba(23, 32, 28, 0.05) 1px, transparent 1px);
      background-size: 44px 44px;
      mask-image: linear-gradient(to bottom, black, transparent 76%);
    }
    main { position: relative; width: min(1440px, calc(100% - 32px)); margin: 0 auto; padding: 34px 0 56px; }
    header {
      display: grid;
      grid-template-columns: minmax(0, 1fr) minmax(320px, 440px);
      gap: 20px;
      align-items: stretch;
      margin-bottom: 18px;
    }
    h1 {
      margin: 0;
      max-width: 860px;
      font-size: clamp(2.5rem, 6vw, 5.6rem);
      letter-spacing: -0.07em;
      line-height: 0.9;
    }
    h2 { margin: 0; font-size: 1.1rem; letter-spacing: -0.02em; }
    p { margin: 0; }
    button, input, select {
      font: inherit;
      border: 1px solid var(--line);
      border-radius: 14px;
      background: rgba(255, 253, 246, 0.92);
      color: var(--ink);
    }
    input, select { min-width: 0; width: 100%; padding: 10px 12px; }
    button {
      cursor: pointer;
      padding: 10px 13px;
      background: var(--accent);
      color: #fffdf6;
      border-color: rgba(37, 91, 73, 0.8);
      transition: transform 160ms ease, opacity 160ms ease, box-shadow 160ms ease;
    }
    button:hover { transform: translateY(-1px); box-shadow: 0 10px 24px rgba(37, 91, 73, 0.18); }
    button:disabled { cursor: wait; opacity: 0.55; transform: none; box-shadow: none; }
    button.secondary { background: rgba(255, 253, 246, 0.86); color: var(--ink); border-color: var(--line); }
    button.ghost { background: transparent; color: var(--accent); border-color: rgba(37, 91, 73, 0.22); }
    .hero, .card, .panel {
      border: 1px solid var(--line);
      background: var(--card);
      box-shadow: var(--shadow);
      border-radius: 28px;
      backdrop-filter: blur(14px);
    }
    .hero { padding: clamp(24px, 4vw, 42px); overflow: hidden; position: relative; }
    .hero::after {
      content: "";
      position: absolute;
      right: -80px;
      bottom: -120px;
      width: 320px;
      height: 320px;
      border-radius: 999px;
      background: rgba(196, 111, 56, 0.16);
    }
    .eyebrow {
      display: inline-flex;
      gap: 8px;
      align-items: center;
      color: var(--accent);
      text-transform: uppercase;
      letter-spacing: 0.12em;
      font-size: 0.78rem;
      font-weight: 700;
      margin-bottom: 18px;
    }
    .sub { color: var(--muted); line-height: 1.6; max-width: 760px; margin-top: 16px; }
    .auth { padding: 18px; display: flex; flex-direction: column; justify-content: space-between; gap: 14px; }
    .auth label, .field label {
      display: block;
      color: var(--muted);
      font-size: 0.74rem;
      text-transform: uppercase;
      letter-spacing: 0.09em;
      margin-bottom: 7px;
    }
    .row { display: flex; gap: 10px; align-items: center; }
    .row.wrap { flex-wrap: wrap; }
    .row > input { flex: 1; }
    .statusLine { color: var(--muted); font-size: 0.9rem; line-height: 1.5; }
    .toolbar, .panel { margin-top: 16px; padding: 16px; }
    .toolbarHead, .panelHead {
      display: flex;
      justify-content: space-between;
      gap: 12px;
      align-items: flex-start;
      margin-bottom: 14px;
    }
    .segmented { display: flex; flex-wrap: wrap; gap: 8px; }
    .segmented button { padding: 8px 12px; border-radius: 999px; }
    .segmented button.active { background: var(--accent-2); border-color: var(--accent-2); color: #fffdf6; }
    .filterGrid { display: grid; grid-template-columns: repeat(6, minmax(0, 1fr)); gap: 10px; }
    .searchField { grid-column: span 2; }
    .metaPills { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 12px; }
    .pill {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 5px 9px;
      border: 1px solid var(--line);
      border-radius: 999px;
      background: rgba(255, 253, 246, 0.68);
      color: var(--muted);
      font-size: 0.78rem;
      white-space: nowrap;
    }
    .pill.good { color: var(--good); border-color: rgba(25, 118, 74, 0.28); background: rgba(25, 118, 74, 0.08); }
    .pill.warn { color: var(--warn); border-color: rgba(182, 106, 19, 0.28); background: rgba(182, 106, 19, 0.08); }
    .pill.bad { color: var(--bad); border-color: rgba(177, 61, 53, 0.28); background: rgba(177, 61, 53, 0.08); }
    .summaryGrid { display: grid; grid-template-columns: repeat(9, minmax(0, 1fr)); gap: 12px; margin-top: 16px; }
    .metric { min-height: 122px; padding: 16px; }
    .metric span { display: block; color: var(--muted); font-size: 0.73rem; text-transform: uppercase; letter-spacing: 0.09em; }
    .metric strong { display: block; margin-top: 10px; font-size: clamp(1.45rem, 2.7vw, 2.7rem); letter-spacing: -0.05em; line-height: 1; }
    .metric small { display: block; margin-top: 10px; color: var(--muted); font-size: 0.82rem; line-height: 1.35; }
    .tableWrap { overflow: auto; border: 1px solid var(--line); border-radius: 22px; background: rgba(255, 253, 246, 0.78); }
    table { width: 100%; min-width: 1080px; border-collapse: collapse; }
    th, td { padding: 12px 13px; border-bottom: 1px solid rgba(36, 52, 45, 0.1); vertical-align: top; text-align: left; }
    th { position: sticky; top: 0; z-index: 2; background: rgba(248, 244, 233, 0.94); color: var(--muted); font-size: 0.73rem; text-transform: uppercase; letter-spacing: 0.08em; }
    th button { padding: 0; border: 0; background: transparent; color: inherit; box-shadow: none; text-transform: inherit; letter-spacing: inherit; }
    tr:last-child td { border-bottom: 0; }
    tr.failedRow { background: rgba(177, 61, 53, 0.055); }
    .primaryCell { display: grid; gap: 4px; }
    .primaryCell strong { font-size: 0.94rem; }
    .primaryCell small, .muted { color: var(--muted); }
    .mono { font-family: var(--mono); font-size: 0.86rem; }
    .goodText { color: var(--good); }
    .warnText { color: var(--warn); }
    .badText { color: var(--bad); }
    .pattern { display: inline-flex; gap: 3px; }
    .dot { width: 8px; height: 18px; border-radius: 999px; background: rgba(102, 115, 108, 0.28); }
    .dot.good { background: var(--good); }
    .dot.bad { background: var(--bad); }
    .pagination { display: flex; justify-content: space-between; align-items: center; gap: 10px; margin-top: 12px; color: var(--muted); font-size: 0.9rem; }
    .pagination .row button { padding: 8px 11px; }
    .error { margin-top: 10px; color: var(--bad); white-space: pre-wrap; font-family: var(--mono); font-size: 0.86rem; }
    .empty { padding: 28px; text-align: center; color: var(--muted); }
    .modalBackdrop {
      position: fixed;
      inset: 0;
      display: none;
      align-items: center;
      justify-content: center;
      padding: 20px;
      background: rgba(23, 32, 28, 0.34);
      z-index: 20;
    }
    .modalBackdrop.open { display: flex; }
    .modal {
      width: min(780px, 100%);
      max-height: min(86vh, 900px);
      overflow: auto;
      border-radius: 28px;
      border: 1px solid var(--line);
      background: #fffdf6;
      box-shadow: 0 28px 100px rgba(23, 32, 28, 0.36);
      padding: 20px;
    }
    .priceGrid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; margin: 14px 0; }
    .priceList { display: grid; gap: 8px; margin-top: 14px; max-height: 260px; overflow: auto; }
    .priceItem { display: flex; justify-content: space-between; gap: 12px; padding: 10px 12px; border: 1px solid var(--line); border-radius: 16px; background: rgba(238, 241, 232, 0.44); }
    .priceItem span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
    .hidden { display: none; }
    @media (max-width: 1180px) {
      header { grid-template-columns: 1fr; }
      .summaryGrid { grid-template-columns: repeat(3, minmax(0, 1fr)); }
      .filterGrid { grid-template-columns: repeat(3, minmax(0, 1fr)); }
      .searchField { grid-column: span 3; }
    }
    @media (max-width: 720px) {
      main { width: min(100% - 20px, 1440px); padding-top: 18px; }
      .hero, .auth, .toolbar, .panel { border-radius: 22px; }
      .summaryGrid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
      .filterGrid, .priceGrid { grid-template-columns: 1fr; }
      .searchField { grid-column: auto; }
      .toolbarHead, .panelHead, .pagination { display: grid; }
      .row { flex-wrap: wrap; }
      .row > button { flex: 1; }
    }
  </style>
</head>
<body>
  <main>
    <header>
      <section class="hero">
        <div class="eyebrow">CLIProxyAPI built-in monitor</div>
        <h1>Request Monitor</h1>
        <p class="sub">Usage analytics now run inside CLIProxyAPI. This panel uses the built-in management endpoints for request history, model prices, cost estimation, account aggregation, filtering, import and export.</p>
      </section>
      <section class="auth card">
        <div>
          <label for="key">Management Key</label>
          <div class="row">
            <input id="key" type="password" autocomplete="current-password" placeholder="Bearer token for /v0/management">
            <button id="saveKey">Save</button>
            <button id="clearKey" class="secondary">Clear</button>
          </div>
        </div>
        <div class="statusLine" id="serviceStatus">Status not loaded.</div>
        <div class="row wrap">
          <button id="refresh">Refresh</button>
          <button id="exportBtn" class="secondary">Export JSONL</button>
          <button id="importBtn" class="secondary">Import JSONL</button>
          <button id="priceBtn" class="ghost">Model prices</button>
          <input id="importFile" class="hidden" type="file" accept=".jsonl,.ndjson,.txt,application/x-ndjson">
        </div>
        <div id="error" class="error"></div>
      </section>
    </header>

    <section class="toolbar panel">
      <div class="toolbarHead">
        <div>
          <h2>Filters</h2>
          <p class="statusLine" id="filterStatus">No data loaded.</p>
        </div>
        <div class="row wrap">
          <select id="autoRefresh" aria-label="Auto refresh">
            <option value="0">Auto refresh off</option>
            <option value="5000">Every 5s</option>
            <option value="10000">Every 10s</option>
            <option value="30000">Every 30s</option>
            <option value="60000">Every 60s</option>
            <option value="300000">Every 5m</option>
          </select>
          <button id="clearFilters" class="secondary">Clear filters</button>
        </div>
      </div>
      <div class="segmented" id="rangeButtons">
        <button data-range="today" class="active">Today</button>
        <button data-range="7d">7 days</button>
        <button data-range="14d">14 days</button>
        <button data-range="30d">30 days</button>
        <button data-range="all">All</button>
      </div>
      <div class="filterGrid" style="margin-top: 12px;">
        <div class="field searchField">
          <label for="search">Search</label>
          <input id="search" type="search" placeholder="account, auth, model, endpoint, provider">
        </div>
        <div class="field">
          <label for="accountFilter">Account</label>
          <select id="accountFilter"></select>
        </div>
        <div class="field">
          <label for="providerFilter">Provider</label>
          <select id="providerFilter"></select>
        </div>
        <div class="field">
          <label for="modelFilter">Model</label>
          <select id="modelFilter"></select>
        </div>
        <div class="field">
          <label for="endpointFilter">Endpoint</label>
          <select id="endpointFilter"></select>
        </div>
        <div class="field">
          <label for="statusFilter">Status</label>
          <select id="statusFilter">
            <option value="all">All statuses</option>
            <option value="success">Success</option>
            <option value="failed">Failed</option>
          </select>
        </div>
      </div>
      <div class="metaPills" id="activeMeta"></div>
    </section>

    <section class="summaryGrid" id="summaryGrid"></section>

    <section class="panel">
      <div class="panelHead">
        <div>
          <h2>Account Overview</h2>
          <p class="statusLine">Grouped by masked source and auth index. Sort columns to find hot accounts, failures, token use, and cost.</p>
        </div>
        <div class="pill" id="accountCount">0 accounts</div>
      </div>
      <div class="tableWrap">
        <table>
          <thead>
            <tr>
              <th>Account</th>
              <th><button data-account-sort="totalCalls">Calls</button></th>
              <th><button data-account-sort="successRate">Success</button></th>
              <th><button data-account-sort="failureCalls">Failures</button></th>
              <th><button data-account-sort="totalTokens">Tokens</button></th>
              <th><button data-account-sort="totalCost">Cost</button></th>
              <th><button data-account-sort="averageLatencyMs">Avg latency</button></th>
              <th><button data-account-sort="lastSeenAt">Latest</button></th>
              <th>Recent</th>
            </tr>
          </thead>
          <tbody id="accountRows"></tbody>
        </table>
      </div>
      <div class="pagination">
        <span id="accountPageInfo">-</span>
        <div class="row">
          <button id="accountPrev" class="secondary">Prev</button>
          <button id="accountNext" class="secondary">Next</button>
        </div>
      </div>
    </section>

    <section class="panel">
      <div class="panelHead">
        <div>
          <h2>Realtime Request Log</h2>
          <p class="statusLine">Latest filtered request rows with per-call tokens, latency, cost, endpoint and auth metadata.</p>
        </div>
        <div class="row">
          <select id="pageSize" aria-label="Rows per page">
            <option value="25">25 rows</option>
            <option value="50" selected>50 rows</option>
            <option value="100">100 rows</option>
            <option value="200">200 rows</option>
            <option value="500">500 rows</option>
          </select>
        </div>
      </div>
      <div class="tableWrap">
        <table>
          <thead>
            <tr>
              <th>Time</th>
              <th>Account</th>
              <th>Provider</th>
              <th>Model</th>
              <th>Endpoint</th>
              <th>Status</th>
              <th>Tokens</th>
              <th>Latency</th>
              <th>Cost</th>
            </tr>
          </thead>
          <tbody id="requestRows"></tbody>
        </table>
      </div>
      <div class="pagination">
        <span id="requestPageInfo">-</span>
        <div class="row">
          <button id="requestPrev" class="secondary">Prev</button>
          <button id="requestNext" class="secondary">Next</button>
        </div>
      </div>
    </section>
  </main>

  <div class="modalBackdrop" id="priceModal">
    <section class="modal">
      <div class="panelHead">
        <div>
          <h2>Model Prices</h2>
          <p class="statusLine">Prices are USD per 1M tokens. Sync uses the built-in LiteLLM model price importer.</p>
        </div>
        <button id="closePrice" class="secondary">Close</button>
      </div>
      <div class="row wrap">
        <button id="syncPrices">Sync visible models</button>
        <span class="statusLine" id="priceStatus">No prices loaded.</span>
      </div>
      <div class="field" style="margin-top: 14px;">
        <label for="priceModel">Model</label>
        <select id="priceModel"></select>
      </div>
      <div class="priceGrid">
        <div class="field">
          <label for="pricePrompt">Input $/1M</label>
          <input id="pricePrompt" type="number" min="0" step="0.0001" placeholder="0">
        </div>
        <div class="field">
          <label for="priceCompletion">Output $/1M</label>
          <input id="priceCompletion" type="number" min="0" step="0.0001" placeholder="0">
        </div>
        <div class="field">
          <label for="priceCache">Cached $/1M</label>
          <input id="priceCache" type="number" min="0" step="0.0001" placeholder="same as input">
        </div>
      </div>
      <div class="row wrap">
        <button id="savePrice">Save price</button>
        <button id="deletePrice" class="secondary">Delete price</button>
      </div>
      <div class="priceList" id="priceList"></div>
    </section>
  </div>

  <script>
    var storageKey = 'cpa_usage_management_key';
    var state = {
      status: null,
      usage: null,
      prices: {},
      rows: [],
      filtered: [],
      accountRows: [],
      range: 'today',
      page: 1,
      pageSize: 50,
      accountPage: 1,
      accountPageSize: 20,
      accountSort: 'lastSeenAt',
      accountSortDir: 'desc',
      timer: null,
      lastRefreshedAt: null
    };

    var els = {
      key: document.getElementById('key'),
      saveKey: document.getElementById('saveKey'),
      clearKey: document.getElementById('clearKey'),
      refresh: document.getElementById('refresh'),
      exportBtn: document.getElementById('exportBtn'),
      importBtn: document.getElementById('importBtn'),
      importFile: document.getElementById('importFile'),
      priceBtn: document.getElementById('priceBtn'),
      error: document.getElementById('error'),
      serviceStatus: document.getElementById('serviceStatus'),
      filterStatus: document.getElementById('filterStatus'),
      rangeButtons: document.getElementById('rangeButtons'),
      autoRefresh: document.getElementById('autoRefresh'),
      clearFilters: document.getElementById('clearFilters'),
      search: document.getElementById('search'),
      accountFilter: document.getElementById('accountFilter'),
      providerFilter: document.getElementById('providerFilter'),
      modelFilter: document.getElementById('modelFilter'),
      endpointFilter: document.getElementById('endpointFilter'),
      statusFilter: document.getElementById('statusFilter'),
      activeMeta: document.getElementById('activeMeta'),
      summaryGrid: document.getElementById('summaryGrid'),
      accountCount: document.getElementById('accountCount'),
      accountRows: document.getElementById('accountRows'),
      accountPageInfo: document.getElementById('accountPageInfo'),
      accountPrev: document.getElementById('accountPrev'),
      accountNext: document.getElementById('accountNext'),
      requestRows: document.getElementById('requestRows'),
      requestPageInfo: document.getElementById('requestPageInfo'),
      requestPrev: document.getElementById('requestPrev'),
      requestNext: document.getElementById('requestNext'),
      pageSize: document.getElementById('pageSize'),
      priceModal: document.getElementById('priceModal'),
      closePrice: document.getElementById('closePrice'),
      syncPrices: document.getElementById('syncPrices'),
      priceStatus: document.getElementById('priceStatus'),
      priceModel: document.getElementById('priceModel'),
      pricePrompt: document.getElementById('pricePrompt'),
      priceCompletion: document.getElementById('priceCompletion'),
      priceCache: document.getElementById('priceCache'),
      savePrice: document.getElementById('savePrice'),
      deletePrice: document.getElementById('deletePrice'),
      priceList: document.getElementById('priceList')
    };

    els.key.value = sessionStorage.getItem(storageKey) || '';
    els.autoRefresh.value = '5000';

    els.saveKey.onclick = function() {
      sessionStorage.setItem(storageKey, els.key.value.trim());
      loadAll();
    };
    els.clearKey.onclick = function() {
      els.key.value = '';
      sessionStorage.removeItem(storageKey);
    };
    els.refresh.onclick = loadAll;
    els.exportBtn.onclick = exportUsage;
    els.importBtn.onclick = function() { els.importFile.click(); };
    els.importFile.onchange = importUsage;
    els.priceBtn.onclick = function() { openPriceModal(); };
    els.closePrice.onclick = function() { els.priceModal.classList.remove('open'); };
    els.syncPrices.onclick = syncModelPrices;
    els.savePrice.onclick = savePrice;
    els.deletePrice.onclick = deletePrice;
    els.priceModel.onchange = loadPriceDraft;
    els.autoRefresh.onchange = resetTimer;
    els.clearFilters.onclick = clearFilters;
    els.search.oninput = function() { state.page = 1; state.accountPage = 1; applyFilters(); };
    els.accountFilter.onchange = onFilterChange;
    els.providerFilter.onchange = onFilterChange;
    els.modelFilter.onchange = onFilterChange;
    els.endpointFilter.onchange = onFilterChange;
    els.statusFilter.onchange = onFilterChange;
    els.pageSize.onchange = function() {
      state.pageSize = Number(els.pageSize.value) || 50;
      state.page = 1;
      renderRequests();
    };
    els.accountPrev.onclick = function() { state.accountPage = Math.max(1, state.accountPage - 1); renderAccounts(); };
    els.accountNext.onclick = function() { state.accountPage += 1; renderAccounts(); };
    els.requestPrev.onclick = function() { state.page = Math.max(1, state.page - 1); renderRequests(); };
    els.requestNext.onclick = function() { state.page += 1; renderRequests(); };
    els.rangeButtons.addEventListener('click', function(event) {
      var btn = event.target.closest('button[data-range]');
      if (!btn) return;
      state.range = btn.getAttribute('data-range') || 'today';
      Array.prototype.forEach.call(els.rangeButtons.querySelectorAll('button'), function(item) {
        item.classList.toggle('active', item === btn);
      });
      state.page = 1;
      state.accountPage = 1;
      applyFilters();
    });
    document.addEventListener('click', function(event) {
      var sortBtn = event.target.closest('button[data-account-sort]');
      if (!sortBtn) return;
      var key = sortBtn.getAttribute('data-account-sort');
      if (state.accountSort === key) {
        state.accountSortDir = state.accountSortDir === 'desc' ? 'asc' : 'desc';
      } else {
        state.accountSort = key;
        state.accountSortDir = 'desc';
      }
      state.accountPage = 1;
      buildAccountRows();
      renderAccounts();
    });

    function onFilterChange() {
      state.page = 1;
      state.accountPage = 1;
      applyFilters();
    }

    function authHeaders() {
      var token = els.key.value.trim() || sessionStorage.getItem(storageKey) || '';
      return token ? { Authorization: 'Bearer ' + token } : {};
    }

    async function fetchJSON(path, options) {
      var init = options || {};
      init.headers = Object.assign({}, authHeaders(), init.headers || {});
      var res = await fetch(path, init);
      if (!res.ok) {
        var text = await res.text();
        throw new Error(path + ' failed: ' + res.status + ' ' + text);
      }
      return res.json();
    }

    async function loadAll() {
      setError('');
      els.refresh.disabled = true;
      try {
        var statusPromise = fetchJSON('/status');
        var usagePromise = fetchJSON('/v0/management/usage');
        var pricesPromise = fetchJSON('/v0/management/model-prices').catch(function() { return { prices: {} }; });
        var result = await Promise.all([statusPromise, usagePromise, pricesPromise]);
        state.status = result[0];
        state.usage = result[1] || {};
        state.prices = result[2].prices || {};
        state.rows = flattenUsage(state.usage);
        state.lastRefreshedAt = new Date();
        refreshFilterOptions();
        applyFilters();
        renderStatus();
        renderPrices();
      } catch (err) {
        setError(err);
      } finally {
        els.refresh.disabled = false;
      }
    }

    function flattenUsage(payload) {
      var out = [];
      var apis = payload && payload.apis ? payload.apis : {};
      Object.keys(apis).forEach(function(endpoint) {
        var apiEntry = apis[endpoint] || {};
        var models = apiEntry.models || {};
        Object.keys(models).forEach(function(model) {
          var details = models[model] && Array.isArray(models[model].details) ? models[model].details : [];
          details.forEach(function(detail, index) {
            var timestampMs = parseTimestampMs(detail.timestamp);
            var tokens = readTokens(detail.tokens || {});
            var method = readString(detail.method) || splitEndpoint(endpoint).method;
            var path = readString(detail.path) || splitEndpoint(endpoint).path;
            var provider = readString(detail.provider) || inferProvider(detail.auth_type || detail.auth_index || detail.source || endpoint);
            var authType = readString(detail.auth_type) || provider;
            var authIndex = readString(detail.auth_index) || '-';
            var source = readString(detail.source) || authIndex || '-';
            var row = {
              id: [detail.request_id || '', endpoint, model, detail.timestamp || '', source, index].join('|'),
              requestId: readString(detail.request_id),
              timestamp: readString(detail.timestamp),
              timestampMs: timestampMs,
              endpoint: endpoint || '-',
              method: method || '-',
              path: path || '-',
              model: model || '-',
              provider: provider || '-',
              authType: authType || '-',
              authIndex: authIndex,
              source: source,
              account: source || authIndex || '-',
              failed: detail.failed === true,
              latencyMs: readNullableNumber(detail.latency_ms),
              tokens: tokens,
              totalCost: calculateCost(model, tokens),
              statsIncluded: detail.failed === true || tokens.total > 0
            };
            row.searchText = [
              row.account,
              row.authIndex,
              row.provider,
              row.authType,
              row.model,
              row.endpoint,
              row.method,
              row.path,
              row.requestId
            ].join(' ').toLowerCase();
            out.push(row);
          });
        });
      });
      out.sort(function(a, b) { return b.timestampMs - a.timestampMs; });
      return out;
    }

    function splitEndpoint(endpoint) {
      var text = readString(endpoint);
      var match = text.match(/^(GET|POST|PUT|PATCH|DELETE|OPTIONS|HEAD)\s+(.+)$/i);
      return match ? { method: match[1].toUpperCase(), path: match[2] } : { method: '', path: text };
    }

    function inferProvider(value) {
      var text = readString(value).toLowerCase();
      var providers = ['codex', 'claude', 'anthropic', 'gemini', 'vertex', 'openai', 'antigravity', 'amp', 'kimi', 'grok', 'qwen'];
      for (var i = 0; i < providers.length; i++) {
        if (text.indexOf(providers[i]) >= 0) return providers[i];
      }
      return '';
    }

    function applyFilters() {
      var startMs = rangeStartMs(state.range);
      var query = els.search.value.trim().toLowerCase();
      var account = els.accountFilter.value || 'all';
      var provider = els.providerFilter.value || 'all';
      var model = els.modelFilter.value || 'all';
      var endpoint = els.endpointFilter.value || 'all';
      var status = els.statusFilter.value || 'all';
      var now = Date.now();
      state.filtered = state.rows.filter(function(row) {
        if (row.timestampMs > now || row.timestampMs < startMs) return false;
        if (account !== 'all' && row.account !== account) return false;
        if (provider !== 'all' && row.provider !== provider) return false;
        if (model !== 'all' && row.model !== model) return false;
        if (endpoint !== 'all' && row.endpoint !== endpoint) return false;
        if (status === 'success' && row.failed) return false;
        if (status === 'failed' && !row.failed) return false;
        if (query && row.searchText.indexOf(query) < 0) return false;
        return true;
      });
      buildAccountRows();
      renderFiltersMeta();
      renderSummary();
      renderAccounts();
      renderRequests();
    }

    function rangeStartMs(range) {
      if (range === 'all') return Number.NEGATIVE_INFINITY;
      var now = new Date();
      now.setHours(0, 0, 0, 0);
      var start = now.getTime();
      if (range === '7d') return start - 6 * 24 * 60 * 60 * 1000;
      if (range === '14d') return start - 13 * 24 * 60 * 60 * 1000;
      if (range === '30d') return start - 29 * 24 * 60 * 60 * 1000;
      return start;
    }

    function refreshFilterOptions() {
      setSelectOptions(els.accountFilter, 'All accounts', uniqueSorted(state.rows.map(function(row) { return row.account; })));
      setSelectOptions(els.providerFilter, 'All providers', uniqueSorted(state.rows.map(function(row) { return row.provider; })));
      setSelectOptions(els.modelFilter, 'All models', uniqueSorted(state.rows.map(function(row) { return row.model; })));
      setSelectOptions(els.endpointFilter, 'All endpoints', uniqueSorted(state.rows.map(function(row) { return row.endpoint; })));
      renderPriceModelOptions();
    }

    function setSelectOptions(select, allLabel, values) {
      var previous = select.value || 'all';
      select.innerHTML = '';
      select.appendChild(new Option(allLabel, 'all'));
      values.forEach(function(value) {
        if (value && value !== '-') select.appendChild(new Option(value, value));
      });
      select.value = values.indexOf(previous) >= 0 ? previous : 'all';
    }

    function renderFiltersMeta() {
      var parts = [
        'Rows: ' + fmt(state.filtered.length),
        'Accounts: ' + fmt(state.accountRows.length),
        'Range: ' + state.range,
        state.lastRefreshedAt ? 'Last refresh: ' + state.lastRefreshedAt.toLocaleTimeString() : ''
      ].filter(Boolean);
      els.filterStatus.textContent = parts.join(' · ');
      els.activeMeta.innerHTML = parts.map(function(part) {
        return '<span class="pill">' + esc(part) + '</span>';
      }).join('');
    }

    function buildAccountRows() {
      var grouped = {};
      state.filtered.forEach(function(row) {
        var key = row.account || row.authIndex || row.source || '-';
        if (!grouped[key]) {
          grouped[key] = {
            account: key,
            providers: {},
            auths: {},
            models: {},
            endpoints: {},
            totalCalls: 0,
            successCalls: 0,
            failureCalls: 0,
            inputTokens: 0,
            outputTokens: 0,
            cachedTokens: 0,
            totalTokens: 0,
            totalCost: 0,
            latencySum: 0,
            latencyCount: 0,
            lastSeenAt: 0,
            recentPattern: []
          };
        }
        var item = grouped[key];
        item.providers[row.provider] = true;
        item.auths[row.authIndex] = true;
        item.models[row.model] = true;
        item.endpoints[row.endpoint] = true;
        item.totalCalls++;
        if (row.failed) item.failureCalls++; else item.successCalls++;
        item.inputTokens += row.tokens.input;
        item.outputTokens += row.tokens.output;
        item.cachedTokens += row.tokens.cached;
        item.totalTokens += row.tokens.total;
        item.totalCost += row.totalCost;
        if (row.latencyMs !== null) {
          item.latencySum += row.latencyMs;
          item.latencyCount++;
        }
        item.lastSeenAt = Math.max(item.lastSeenAt, row.timestampMs);
        item.recentPattern.push(!row.failed);
        if (item.recentPattern.length > 10) item.recentPattern.shift();
      });
      state.accountRows = Object.keys(grouped).map(function(key) {
        var item = grouped[key];
        item.successRate = item.totalCalls ? item.successCalls / item.totalCalls : 1;
        item.averageLatencyMs = item.latencyCount ? item.latencySum / item.latencyCount : null;
        item.providerText = joinLimited(Object.keys(item.providers), 2);
        item.authText = joinLimited(Object.keys(item.auths), 2);
        item.modelText = joinLimited(Object.keys(item.models), 3);
        item.endpointText = joinLimited(Object.keys(item.endpoints), 2);
        return item;
      });
      var dir = state.accountSortDir === 'desc' ? -1 : 1;
      state.accountRows.sort(function(a, b) {
        var av = comparable(a[state.accountSort]);
        var bv = comparable(b[state.accountSort]);
        if (av < bv) return -1 * dir;
        if (av > bv) return 1 * dir;
        return b.lastSeenAt - a.lastSeenAt;
      });
    }

    function renderSummary() {
      var rows = state.filtered.filter(function(row) { return row.statsIncluded; });
      var total = rows.length;
      var failures = rows.filter(function(row) { return row.failed; }).length;
      var success = Math.max(total - failures, 0);
      var input = sum(rows, function(row) { return row.tokens.input; });
      var output = sum(rows, function(row) { return row.tokens.output; });
      var reasoning = sum(rows, function(row) { return row.tokens.reasoning; });
      var cached = sum(rows, function(row) { return row.tokens.cached; });
      var tokens = sum(rows, function(row) { return row.tokens.total; });
      var cost = sum(rows, function(row) { return row.totalCost; });
      var latencyRows = rows.filter(function(row) { return row.latencyMs !== null; });
      var avgLatency = latencyRows.length ? sum(latencyRows, function(row) { return row.latencyMs; }) / latencyRows.length : null;
      var recentStart = Date.now() - 30 * 60 * 1000;
      var recentRows = rows.filter(function(row) { return row.timestampMs >= recentStart; });
      var cards = [
        ['Total calls', fmt(total), state.accountRows.length + ' accounts'],
        ['Success', fmt(success), pct(total ? success / total : 1), 'good'],
        ['Failure', fmt(failures), failures + ' failed rows', failures ? 'bad' : 'good'],
        ['Success rate', pct(total ? success / total : 1), dur(avgLatency) + ' avg latency', total && success / total < 0.85 ? 'bad' : total && success / total < 0.95 ? 'warn' : 'good'],
        ['Total tokens', fmt(tokens), 'Reasoning ' + fmt(reasoning)],
        ['Input tokens', fmt(input), pct(tokens ? input / tokens : 0) + ' of token mix'],
        ['Output tokens', fmt(output), pct(tokens ? output / tokens : 0) + ' of token mix'],
        ['Cached tokens', fmt(cached), pct(input ? cached / input : 0) + ' of input'],
        ['Estimated cost', money(cost), Object.keys(state.prices).length ? 'Using saved prices' : 'No model prices', Object.keys(state.prices).length ? '' : 'warn']
      ];
      els.summaryGrid.innerHTML = cards.map(function(card) {
        return '<div class="metric card">' +
          '<span>' + esc(card[0]) + '</span>' +
          '<strong class="' + toneClass(card[3]) + '">' + esc(card[1]) + '</strong>' +
          '<small>' + esc(card[2]) + '</small>' +
          '</div>';
      }).join('');
      els.activeMeta.innerHTML += '<span class="pill">RPM 30m: ' + esc((recentRows.length / 30).toFixed(2)) + '</span>';
    }

    function renderStatus() {
      var s = state.status || {};
      var c = s.collector || {};
      var pieces = [
        'collector=' + (c.collector || '-'),
        'transport=' + (c.transport || '-'),
        'events=' + fmt(s.events),
        'deadLetters=' + fmt(s.deadLetters),
        'db=' + (s.dbPath || '-')
      ];
      if (c.lastError) pieces.push('lastError=' + c.lastError);
      els.serviceStatus.textContent = pieces.join(', ');
    }

    function renderAccounts() {
      var totalPages = Math.max(1, Math.ceil(state.accountRows.length / state.accountPageSize));
      state.accountPage = Math.min(Math.max(1, state.accountPage), totalPages);
      var start = (state.accountPage - 1) * state.accountPageSize;
      var items = state.accountRows.slice(start, start + state.accountPageSize);
      els.accountCount.textContent = fmt(state.accountRows.length) + ' accounts';
      els.accountRows.innerHTML = items.map(function(row) {
        var successClass = row.successRate >= 0.95 ? 'goodText' : row.successRate >= 0.85 ? 'warnText' : 'badText';
        return '<tr>' +
          '<td><div class="primaryCell"><strong>' + esc(row.account) + '</strong><small>' + esc(row.providerText || '-') + ' · ' + esc(row.authText || '-') + '</small></div></td>' +
          '<td>' + fmt(row.totalCalls) + '</td>' +
          '<td class="' + successClass + '">' + pct(row.successRate) + '</td>' +
          '<td class="' + (row.failureCalls ? 'badText' : 'goodText') + '">' + fmt(row.failureCalls) + '</td>' +
          '<td><div class="primaryCell"><span>' + fmt(row.totalTokens) + '</span><small>I ' + fmt(row.inputTokens) + ' · O ' + fmt(row.outputTokens) + ' · C ' + fmt(row.cachedTokens) + '</small></div></td>' +
          '<td>' + money(row.totalCost) + '</td>' +
          '<td>' + dur(row.averageLatencyMs) + '</td>' +
          '<td>' + dateText(row.lastSeenAt) + '</td>' +
          '<td>' + renderPattern(row.recentPattern) + '</td>' +
          '</tr>';
      }).join('') || '<tr><td colspan="9"><div class="empty">No account data for current filters.</div></td></tr>';
      els.accountPageInfo.textContent = pageInfo(state.accountRows.length, state.accountPage, state.accountPageSize);
      els.accountPrev.disabled = state.accountPage <= 1;
      els.accountNext.disabled = state.accountPage >= totalPages;
    }

    function renderRequests() {
      var totalPages = Math.max(1, Math.ceil(state.filtered.length / state.pageSize));
      state.page = Math.min(Math.max(1, state.page), totalPages);
      var start = (state.page - 1) * state.pageSize;
      var items = state.filtered.slice(start, start + state.pageSize);
      els.requestRows.innerHTML = items.map(function(row) {
        return '<tr class="' + (row.failed ? 'failedRow' : '') + '">' +
          '<td><div class="primaryCell"><span>' + esc(dateText(row.timestampMs)) + '</span><small class="mono">' + esc(row.requestId || '-') + '</small></div></td>' +
          '<td><div class="primaryCell"><strong>' + esc(row.account) + '</strong><small class="mono">' + esc(row.authIndex) + '</small></div></td>' +
          '<td>' + esc(row.provider) + '</td>' +
          '<td class="mono">' + esc(row.model) + '</td>' +
          '<td><div class="primaryCell"><span class="mono">' + esc(row.method + ' ' + row.path) + '</span><small>' + esc(row.endpoint) + '</small></div></td>' +
          '<td><span class="pill ' + (row.failed ? 'bad' : 'good') + '">' + (row.failed ? 'failed' : 'success') + '</span></td>' +
          '<td><div class="primaryCell"><span>' + fmt(row.tokens.total) + '</span><small>I ' + fmt(row.tokens.input) + ' · O ' + fmt(row.tokens.output) + ' · R ' + fmt(row.tokens.reasoning) + ' · C ' + fmt(row.tokens.cached) + '</small></div></td>' +
          '<td>' + dur(row.latencyMs) + '</td>' +
          '<td>' + money(row.totalCost) + '</td>' +
          '</tr>';
      }).join('') || '<tr><td colspan="9"><div class="empty">No request rows for current filters.</div></td></tr>';
      els.requestPageInfo.textContent = pageInfo(state.filtered.length, state.page, state.pageSize);
      els.requestPrev.disabled = state.page <= 1;
      els.requestNext.disabled = state.page >= totalPages;
    }

    async function exportUsage() {
      try {
        var res = await fetch('/v0/management/usage/export', { headers: authHeaders() });
        if (!res.ok) throw new Error(await res.text());
        var blob = await res.blob();
        var url = URL.createObjectURL(blob);
        var a = document.createElement('a');
        a.href = url;
        a.download = 'usage-events.jsonl';
        a.click();
        URL.revokeObjectURL(url);
      } catch (err) {
        setError(err);
      }
    }

    async function importUsage() {
      var file = els.importFile.files && els.importFile.files[0];
      if (!file) return;
      try {
        var text = await file.text();
        var res = await fetch('/v0/management/usage/import', {
          method: 'POST',
          headers: Object.assign({}, authHeaders(), { 'Content-Type': 'application/x-ndjson' }),
          body: text
        });
        if (!res.ok) throw new Error(await res.text());
        var result = await res.json();
        setError('Imported ' + fmt(result.added) + ', skipped ' + fmt(result.skipped) + ', failed ' + fmt(result.failed) + '.');
        await loadAll();
      } catch (err) {
        setError(err);
      } finally {
        els.importFile.value = '';
      }
    }

    function openPriceModal() {
      els.priceModal.classList.add('open');
      renderPriceModelOptions();
      loadPriceDraft();
      renderPrices();
    }

    function renderPriceModelOptions() {
      var models = uniqueSorted(state.rows.map(function(row) { return row.model; }).concat(Object.keys(state.prices)));
      els.priceModel.innerHTML = '';
      if (!models.length) {
        els.priceModel.appendChild(new Option('No models loaded', ''));
        return;
      }
      models.forEach(function(model) { els.priceModel.appendChild(new Option(model, model)); });
    }

    function loadPriceDraft() {
      var model = els.priceModel.value;
      var price = state.prices[model] || {};
      els.pricePrompt.value = price.prompt === undefined ? '' : String(price.prompt);
      els.priceCompletion.value = price.completion === undefined ? '' : String(price.completion);
      els.priceCache.value = price.cache === undefined ? '' : String(price.cache);
    }

    async function savePrice() {
      var model = els.priceModel.value;
      if (!model) return;
      var next = Object.assign({}, state.prices);
      var prompt = numberOrZero(els.pricePrompt.value);
      var completion = numberOrZero(els.priceCompletion.value);
      next[model] = {
        prompt: prompt,
        completion: completion,
        cache: els.priceCache.value.trim() === '' ? prompt : numberOrZero(els.priceCache.value)
      };
      await savePrices(next);
    }

    async function deletePrice() {
      var model = els.priceModel.value;
      if (!model) return;
      var next = Object.assign({}, state.prices);
      delete next[model];
      await savePrices(next);
    }

    async function savePrices(next) {
      try {
        var response = await fetchJSON('/v0/management/model-prices', {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ prices: next })
        });
        state.prices = response.prices || next;
        state.rows.forEach(function(row) { row.totalCost = calculateCost(row.model, row.tokens); });
        applyFilters();
        renderPrices();
        loadPriceDraft();
      } catch (err) {
        setError(err);
      }
    }

    async function syncModelPrices() {
      try {
        els.syncPrices.disabled = true;
        var models = uniqueSorted(state.filtered.map(function(row) { return row.model; }));
        var response = await fetchJSON('/v0/management/model-prices/sync', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ models: models })
        });
        state.prices = response.prices || {};
        state.rows.forEach(function(row) { row.totalCost = calculateCost(row.model, row.tokens); });
        applyFilters();
        renderPrices();
        loadPriceDraft();
        els.priceStatus.textContent = 'Synced ' + fmt(response.imported) + ' prices, skipped ' + fmt(response.skipped) + '.';
      } catch (err) {
        setError(err);
      } finally {
        els.syncPrices.disabled = false;
      }
    }

    function renderPrices() {
      var entries = Object.keys(state.prices).sort().map(function(model) {
        return [model, state.prices[model]];
      });
      els.priceStatus.textContent = entries.length ? entries.length + ' prices loaded.' : 'No prices loaded.';
      els.priceList.innerHTML = entries.slice(0, 120).map(function(entry) {
        var p = entry[1] || {};
        return '<div class="priceItem"><span class="mono">' + esc(entry[0]) + '</span><span>' +
          moneyPerM(p.prompt) + ' / ' + moneyPerM(p.completion) + ' / cache ' + moneyPerM(p.cache) +
          '</span></div>';
      }).join('') || '<div class="empty">No model prices saved yet.</div>';
    }

    function calculateCost(model, tokens) {
      var p = state.prices[model];
      if (!p) return 0;
      return ((tokens.input || 0) * numberOrZero(p.prompt) +
        (tokens.output || 0) * numberOrZero(p.completion) +
        (tokens.cached || 0) * numberOrZero(p.cache)) / 1000000;
    }

    function readTokens(raw) {
      var input = numberOrZero(raw.input_tokens !== undefined ? raw.input_tokens : raw.inputTokens);
      var output = numberOrZero(raw.output_tokens !== undefined ? raw.output_tokens : raw.outputTokens);
      var reasoning = numberOrZero(raw.reasoning_tokens !== undefined ? raw.reasoning_tokens : raw.reasoningTokens);
      var cached = numberOrZero(raw.cached_tokens !== undefined ? raw.cached_tokens : raw.cachedTokens);
      var cache = numberOrZero(raw.cache_tokens !== undefined ? raw.cache_tokens : raw.cacheTokens);
      var total = numberOrZero(raw.total_tokens !== undefined ? raw.total_tokens : raw.totalTokens);
      if (!total) total = input + output + reasoning + Math.max(cached, cache);
      return { input: input, output: output, reasoning: reasoning, cached: cached, cache: cache, total: total };
    }

    function clearFilters() {
      els.search.value = '';
      els.accountFilter.value = 'all';
      els.providerFilter.value = 'all';
      els.modelFilter.value = 'all';
      els.endpointFilter.value = 'all';
      els.statusFilter.value = 'all';
      state.page = 1;
      state.accountPage = 1;
      applyFilters();
    }

    function resetTimer() {
      if (state.timer) clearInterval(state.timer);
      state.timer = null;
      var ms = Number(els.autoRefresh.value);
      if (ms > 0) {
        state.timer = setInterval(function() {
          loadAll().catch(function() {});
        }, ms);
      }
    }

    function setError(err) {
      if (!err) {
        els.error.textContent = '';
        return;
      }
      els.error.textContent = err.message ? err.message : String(err);
    }

    function pageInfo(total, page, pageSize) {
      if (!total) return '0 rows';
      var start = (page - 1) * pageSize + 1;
      var end = Math.min(start + pageSize - 1, total);
      return fmt(start) + '-' + fmt(end) + ' of ' + fmt(total);
    }

    function renderPattern(pattern) {
      var data = pattern && pattern.length ? pattern : [];
      return '<span class="pattern">' + data.map(function(ok) {
        return '<span class="dot ' + (ok ? 'good' : 'bad') + '"></span>';
      }).join('') + '</span>';
    }

    function uniqueSorted(values) {
      var seen = {};
      values.forEach(function(value) {
        var text = readString(value);
        if (text) seen[text] = true;
      });
      return Object.keys(seen).sort(function(a, b) { return a.localeCompare(b); });
    }

    function joinLimited(values, limit) {
      var clean = uniqueSorted(values).filter(function(value) { return value !== '-'; });
      if (clean.length <= limit) return clean.join(', ');
      return clean.slice(0, limit).join(', ') + ' +' + (clean.length - limit);
    }

    function sum(rows, fn) {
      return rows.reduce(function(total, row) { return total + numberOrZero(fn(row)); }, 0);
    }

    function comparable(value) {
      if (value === null || value === undefined) return -1;
      if (typeof value === 'number') return value;
      return String(value).toLowerCase();
    }

    function readString(value) {
      if (value === null || value === undefined) return '';
      return String(value).trim();
    }

    function numberOrZero(value) {
      var n = Number(value);
      return Number.isFinite(n) && n >= 0 ? n : 0;
    }

    function readNullableNumber(value) {
      if (value === null || value === undefined || value === '') return null;
      var n = Number(value);
      return Number.isFinite(n) && n >= 0 ? n : null;
    }

    function parseTimestampMs(value) {
      var parsed = Date.parse(value || '');
      return Number.isFinite(parsed) ? parsed : 0;
    }

    function fmt(value) {
      if (value === null || value === undefined || value === '') return '-';
      var n = Number(value);
      return Number.isFinite(n) ? n.toLocaleString() : String(value);
    }

    function pct(value) {
      var n = Number(value);
      return Number.isFinite(n) ? (n * 100).toFixed(1) + '%' : '-';
    }

    function money(value) {
      var n = Number(value);
      if (!Number.isFinite(n)) return '-';
      if (Math.abs(n) < 0.01) return '$' + n.toFixed(4);
      return '$' + n.toFixed(2);
    }

    function moneyPerM(value) {
      return '$' + numberOrZero(value).toFixed(4) + '/1M';
    }

    function dur(value) {
      if (value === null || value === undefined || !Number.isFinite(Number(value))) return '-';
      var ms = Number(value);
      if (ms < 1000) return Math.round(ms) + ' ms';
      if (ms < 60000) return (ms / 1000).toFixed(1) + ' s';
      return (ms / 60000).toFixed(1) + ' min';
    }

    function dateText(ms) {
      if (!ms) return '-';
      return new Date(ms).toLocaleString();
    }

    function toneClass(tone) {
      if (tone === 'good') return 'goodText';
      if (tone === 'warn') return 'warnText';
      if (tone === 'bad') return 'badText';
      return '';
    }

    function esc(value) {
      return String(value === undefined || value === null ? '' : value).replace(/[&<>"']/g, function(ch) {
        if (ch === '&') return '&amp;';
        if (ch === '<') return '&lt;';
        if (ch === '>') return '&gt;';
        if (ch === '"') return '&quot;';
        return '&#39;';
      });
    }

    resetTimer();
    loadAll();
  </script>
</body>
</html>`
