// The page. Fetches the published copies of evals/results/*.jsonl, summarises
// each entry with stats.js, and renders. It computes no number stats.js does
// not, and writes nothing anywhere.

import { wilson, fmtInterval, summarise, compare, BASELINE_UNDEFINED } from "./stats.js";

// ── tiny helpers ──────────────────────────────────────────────────

const $ = (sel) => document.querySelector(sel);

function esc(s) {
  return String(s ?? "").replace(/[&<>"']/g, (c) =>
    ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" })[c],
  );
}

const fmtDate = (iso) => {
  const d = new Date(iso);
  return Number.isNaN(+d) ? String(iso).slice(0, 10) : d.toISOString().slice(0, 10);
};

function fmtDuration(ms) {
  if (!ms) return "—";
  if (ms < 1000) return `${Math.round(ms)}ms`;
  if (ms < 60_000) return `${(ms / 1000).toFixed(1)}s`;
  const m = Math.floor(ms / 60_000);
  const s = Math.round((ms % 60_000) / 1000);
  return m >= 60 ? `${Math.floor(m / 60)}h ${m % 60}m` : `${m}m ${s}s`;
}

// thousands keeps a token count readable without losing the small ones, on the
// same rule the terminal report follows.
const thousands = (n) => (n < 10_000 ? String(n) : `${(n / 1000).toFixed(1)}k`);

// ── the one chart form on this page ───────────────────────────────
//
// A single-series interval strip: the band spans the interval, the rule marks
// the point estimate, and the fraction is always printed beside it. One hue,
// so no legend — the column header names the series.
function strip(iv, title) {
  if (!iv) return `<span class="iv-none">not measured</span>`;
  const left = (iv.low * 100).toFixed(2);
  const width = Math.max(0.8, (iv.high - iv.low) * 100).toFixed(2);
  const point = (iv.rate * 100).toFixed(2);
  const label = `${fmtInterval(iv)} — the interval is the range of true rates consistent with ${iv.samples} sample(s)`;
  return `<span class="iv" title="${esc(title ? `${title}. ${label}` : label)}">
    <span class="iv-text">${esc(fmtInterval(iv))}</span>
    <span class="iv-track" role="img" aria-label="${esc(label)}">
      <span class="iv-band" style="left:${left}%;width:${width}%"></span>
      <span class="iv-point" style="left:calc(${point}% - 1.5px)"></span>
    </span>
  </span>`;
}

function blank(reason) {
  return `<span class="blank-cell" title="${esc(reason)}">not defined</span>`;
}

// The resolution badge beside n already says "smoke" and "unstated", so the
// Flags cell carries only what it does not: an unpinnable tree, unknown inputs.
const RESOLUTION_TAGS = new Set(["smoke", "unstated"]);

function flagsFor(s) {
  const cls = { dirty: "flag-dirty", "no digest": "flag-unstated" };
  return s.caveats
    .filter((c) => !RESOLUTION_TAGS.has(c.tag))
    .map((c) => `<span class="flag ${cls[c.tag] || ""}" title="${esc(c.why)}">${esc(c.tag)}</span>`)
    .join("");
}

// A resolution that cannot be cited is marked where it is stated, and carries
// the reason it cannot.
function resolutionBadge(s) {
  const tag = s.resolution || "unstated";
  const caveat = s.caveats.find((c) => RESOLUTION_TAGS.has(c.tag));
  return `<span class="res res-${esc(tag)}${caveat ? " res-flagged" : ""}"${
    caveat ? ` title="${esc(caveat.why)}"` : ""
  }>${esc(tag)}</span>`;
}

// ── state ─────────────────────────────────────────────────────────

const state = {
  rows: [],
  cases: {},
  source: {},
  filters: { case: "", model: "", arm: "", framework: "", tightness: "", resolution: "", citable: false },
  sort: { key: "at", dir: "desc" },
};

// ── load ──────────────────────────────────────────────────────────

async function load() {
  const manifest = await (await fetch("data/manifest.json")).json();
  state.source = manifest.source || {};

  const rows = [];
  for (const c of manifest.cases) {
    if (c.case) state.cases[c.id] = c.case;
    const text = await (await fetch(c.file)).text();
    for (const line of text.split("\n")) {
      if (!line.trim()) continue;
      try {
        rows.push(summarise(JSON.parse(line)));
      } catch (e) {
        // An entry that will not parse is reported, never repaired.
        rows.push({ unreadable: true, case: c.id, error: String(e), caveats: [], comparable: false });
      }
    }
  }
  state.rows = rows.filter((r) => !r.unreadable);

  const unreadable = rows.filter((r) => r.unreadable).length + (manifest.unreadable?.length || 0);

  $("#built").textContent =
    `${state.rows.length} runs across ${manifest.cases.length} cases · ` +
    `${state.rows.reduce((a, r) => a + r.samples, 0)} samples · ` +
    `built from ${state.source.commit || "an unknown commit"}` +
    (unreadable ? ` · ${unreadable} unreadable line(s)` : "");

  $("#source-note").textContent = state.source.commit
    ? `Built from commit ${state.source.commit}.`
    : "";
}

// ── controls ──────────────────────────────────────────────────────

const AXES = [
  ["case", "Case", (r) => r.case],
  ["model", "Model", (r) => r.model.ID],
  ["arm", "Arm", (r) => r.arm],
  ["framework", "Framework", (r) => r.framework],
  ["tightness", "Tightness", (r) => r.tightness],
  ["resolution", "Resolution", (r) => r.resolution || "(unstated)"],
];

function renderControls() {
  const host = $("#controls");
  host.innerHTML = AXES.map(([key, label, get]) => {
    const values = [...new Set(state.rows.map(get).filter(Boolean))].sort();
    return `<span class="ctl">
      <label for="f-${key}">${esc(label)}</label>
      <select id="f-${key}" data-filter="${key}">
        <option value="">All</option>
        ${values.map((v) => `<option value="${esc(v)}">${esc(v)}</option>`).join("")}
      </select>
    </span>`;
  }).join("") +
    `<span class="ctl ctl-check">
       <input type="checkbox" id="f-citable">
       <label for="f-citable">Only runs that can be cited</label>
     </span>
     <button class="btn-reset" id="reset">Reset</button>`;

  host.querySelectorAll("select[data-filter]").forEach((sel) => {
    sel.addEventListener("change", () => {
      state.filters[sel.dataset.filter] = sel.value;
      render();
    });
  });
  $("#f-citable").addEventListener("change", (e) => {
    state.filters.citable = e.target.checked;
    render();
  });
  $("#reset").addEventListener("click", () => {
    for (const k of Object.keys(state.filters)) state.filters[k] = k === "citable" ? false : "";
    host.querySelectorAll("select").forEach((s) => (s.value = ""));
    $("#f-citable").checked = false;
    render();
  });
}

function visible() {
  const f = state.filters;
  return state.rows.filter((r) => {
    if (f.case && r.case !== f.case) return false;
    if (f.model && r.model.ID !== f.model) return false;
    if (f.arm && r.arm !== f.arm) return false;
    if (f.framework && r.framework !== f.framework) return false;
    if (f.tightness && r.tightness !== f.tightness) return false;
    if (f.resolution && (r.resolution || "(unstated)") !== f.resolution) return false;
    if (f.citable && !r.comparable) return false;
    return true;
  });
}

// ── sorting ───────────────────────────────────────────────────────
//
// Rate columns sort on the point estimate and push "not measured" to the
// bottom in both directions — a run that did not measure something is not the
// worst at it.
const SORT_KEYS = {
  at: (r) => r.at,
  case: (r) => r.case,
  model: (r) => `${r.model.ID} ${r.model.Engine}`,
  arm: (r) => r.arm,
  samples: (r) => r.samples,
  validity: (r) => r.validity?.rate,
  behavior: (r) => r.behavior?.workingRate?.rate,
  fill: (r) => r.fill?.interval?.rate,
};

function sorted(rows) {
  const get = SORT_KEYS[state.sort.key] || SORT_KEYS.at;
  const dir = state.sort.dir === "asc" ? 1 : -1;
  return [...rows].sort((a, b) => {
    const x = get(a), y = get(b);
    const xm = x === undefined || x === null, ym = y === undefined || y === null;
    if (xm && ym) return 0;
    if (xm) return 1;           // missing always last
    if (ym) return -1;
    if (x === y) return b.at.localeCompare(a.at);
    return (x > y ? 1 : -1) * dir;
  });
}

function wireSorting() {
  document.querySelectorAll("th.sortable").forEach((th) => {
    th.addEventListener("click", () => {
      const key = th.dataset.sort;
      if (state.sort.key === key) state.sort.dir = state.sort.dir === "asc" ? "desc" : "asc";
      else state.sort = { key, dir: key === "case" || key === "model" || key === "arm" ? "asc" : "desc" };
      render();
    });
  });
}

// ── table ─────────────────────────────────────────────────────────

function renderTable(rows) {
  document.querySelectorAll("th.sortable").forEach((th) => {
    th.classList.remove("sorted-asc", "sorted-desc");
    if (th.dataset.sort === state.sort.key) th.classList.add(`sorted-${state.sort.dir}`);
  });

  const body = $("#runs-body");
  $("#runs-empty").hidden = rows.length > 0;

  body.innerHTML = sorted(rows)
    .map((s, i) => {
      const behaviour = s.behavior
        ? strip(s.behavior.workingRate, "samples whose application built, booted and passed every assertion")
        : `<span class="iv-none">not measured</span>`;

      const fill = s.baseline
        ? blank(BASELINE_UNDEFINED)
        : s.fill
          ? strip(s.fill.interval, "planted anchors a selection accounted for")
          : `<span class="iv-none">not recorded</span>`;

      return `<tr data-i="${i}" class="${s.comparable ? "" : "not-citable"}" tabindex="0">
        <td class="cell-date keep">${esc(fmtDate(s.at))}</td>
        <td class="cell-case">${esc(s.case)}<span class="stack">${esc([s.language, s.framework, s.tightness].filter(Boolean).join(" · "))}</span></td>
        <td class="cell-model">${esc(s.model.ID || "—")}<span class="engine">${esc([s.model.Engine, s.model.Quant].filter(Boolean).join(" · "))}</span></td>
        <td><span class="armbadge armbadge-${esc(s.arm)}">${esc(s.arm)}</span></td>
        <td class="num">${s.samples} ${resolutionBadge(s)}</td>
        <td>${strip(s.validity, "samples the pipeline accepted")}${
          s.failed ? `<span class="iv-note">${s.failed} of ${s.samples} never reached the model</span>` : ""
        }</td>
        <td>${behaviour}</td>
        <td>${fill}</td>
        <td class="keep">${flagsFor(s)}</td>
      </tr>`;
    })
    .join("");

  const ordered = sorted(rows);
  body.querySelectorAll("tr").forEach((tr) => {
    const open = () => openDetail(ordered[Number(tr.dataset.i)]);
    tr.addEventListener("click", open);
    tr.addEventListener("keydown", (e) => {
      if (e.key === "Enter" || e.key === " ") { e.preventDefault(); open(); }
    });
  });
}

function renderSummary(rows) {
  const citable = rows.filter((r) => r.comparable).length;
  const excluded = rows.length - citable;
  const samples = rows.reduce((a, r) => a + r.samples, 0);
  const fixtures = new Set(rows.map((r) => r.fixture || "(none)")).size;

  $("#summary-bar").innerHTML =
    `Showing <strong>${rows.length}</strong> of <strong>${state.rows.length}</strong> runs, ` +
    `<strong>${samples}</strong> samples in total. ` +
    `<strong>${citable}</strong> can be cited as measurements` +
    (excluded
      ? `; <span class="excluded">${excluded} are shown but excluded from every comparison — hover a flag for the reason.</span>`
      : ".") +
    (fixtures > 1
      ? ` <span class="excluded">These runs span ${fixtures} different fixture digests, so nothing here is averaged across them.</span>`
      : "");
}

// ── detail ────────────────────────────────────────────────────────

function statList(pairs) {
  return `<div class="dgrid">${pairs
    .filter(([, v]) => v !== null && v !== undefined && v !== "")
    .map(([k, v]) => `<span class="stat"><span class="k">${esc(k)}</span><span class="v">${esc(v)}</span></span>`)
    .join("")}</div>`;
}

function countTable(head, list, total) {
  if (!list || !list.length) return "";
  return `<table class="mini"><thead><tr><th>${esc(head)}</th><th class="num">samples</th></tr></thead><tbody>
    ${list.map((d) => `<tr><td>${esc(d.name)}</td><td class="num">${d.n}${total ? ` / ${total}` : ""}</td></tr>`).join("")}
  </tbody></table>`;
}

function group(title, inner) {
  return inner ? `<div class="dgroup"><h4>${esc(title)}</h4>${inner}</div>` : "";
}

function openDetail(s) {
  const meta = state.cases[s.case];

  const caveats = s.caveats.length
    ? s.caveats
        .map((c) => `<p class="caveat"><strong>${esc(c.tag)}</strong> — ${esc(c.why)}</p>`)
        .join("")
    : "";

  const outcome = group(
    "Outcome",
    `<div class="outcomes">
      <span class="outcome"><span class="dot dot-good"></span><span class="n">${s.valid}</span> valid</span>
      <span class="outcome"><span class="dot dot-serious"></span><span class="n">${s.invalid}</span> invalid (answered, rejected)</span>
      <span class="outcome"><span class="dot dot-muted"></span><span class="n">${s.failed}</span> failed (never reached the model)</span>
    </div>
    <div class="cmp-row"><span class="cmp-label">valid of answered</span>${strip(s.validity)}</div>
    ${s.firstCall ? `<div class="cmp-row"><span class="cmp-label">valid on first call</span>${strip(s.firstCall)}</div>` : ""}
    ${s.failed ? `<p class="note">The ${s.failed} failed sample(s) are outside the denominator: an unreachable endpoint is not a model that chose badly.</p>` : ""}
    ${countTable("rejected under rule", s.rules)}
    ${s.details.length ? `<details class="more"><summary>${s.details.length} distinct reason(s) recorded</summary><pre class="inv">${esc(s.details.join("\n"))}</pre></details>` : ""}`,
  );

  const provenance = group(
    "How this run was drawn",
    statList([
      ["commit", `${s.commit}${s.clean ? "" : " (dirty)"}`],
      ["fixture digest", s.fixture || "not recorded"],
      ["resolution", s.resolution || "unstated"],
      ["samples", s.samples],
      ["concurrency", s.concurrency],
      ["retry budget", s.retries],
      ["endpoint", s.endpoint],
      ["arm", s.arm],
      ["stack", [s.language, s.framework, s.tightness].filter(Boolean).join(" · ")],
    ]) + (s.dirty.length ? `<p class="note">Uncommitted at run time: ${esc(s.dirty.join(", "))}</p>` : ""),
  );

  const cost = group(
    "What it cost",
    statList([
      ["wall clock", fmtDuration(s.wallMs)],
      ["per sample", s.timing ? `${fmtDuration(s.timing.fastest)} / ${fmtDuration(s.timing.mean)} / ${fmtDuration(s.timing.slowest)}` : null],
      ["model calls", s.calls ? `${s.calls.total} over ${s.calls.over} sample(s), mean ${s.calls.mean.toFixed(2)}` : null],
      ["prompt tokens", s.tokens ? thousands(s.tokens.prompt) : null],
      ["completion tokens", s.tokens ? thousands(s.tokens.completion) : null],
    ]) + (s.timing ? `<p class="note">Per sample: fastest / mean / slowest.</p>` : ""),
  );

  const behaviour = s.behavior
    ? group(
        "Behaviour — the application was built, booted and asserted against",
        `<div class="outcomes">
          <span class="outcome"><span class="dot dot-good"></span><span class="n">${s.behavior.working}</span> working</span>
          <span class="outcome"><span class="dot dot-serious"></span><span class="n">${s.behavior.disagreed}</span> disagreed (booted, failed an assertion)</span>
          <span class="outcome"><span class="dot dot-critical"></span><span class="n">${s.behavior.broke}</span> broke (never reached the assertions)</span>
        </div>
        <div class="cmp-row"><span class="cmp-label">working</span>${strip(s.behavior.workingRate)}</div>
        <p class="note">${s.behavior.passed} of ${s.behavior.checks} assertions held across ${s.behavior.measured} applied sample(s). That total is a bare ratio, not a rate over samples, so it carries no interval.</p>
        ${countTable("assertion that did not hold", s.behavior.failures, s.behavior.measured)}
        ${countTable("phase that died", s.behavior.phases)}`,
      )
    : group(
        "Behaviour",
        `<p class="note">Not measured on this run. Behaviour is off by default — one sample costs a scaffold, a dependency install, a database and a boot. This is not a run whose application failed; it is a run where nothing was booted.</p>`,
      );

  const signals = s.baseline
    ? group("Form", `<p class="note">${esc(BASELINE_UNDEFINED)}</p>`)
    : group(
        "Form — cheap checks over what was rendered",
        [
          s.fill
            ? `<div class="cmp-row"><span class="cmp-label">anchors filled</span>${strip(s.fill.interval)}</div>
               <p class="note">${s.fill.numer} of ${s.fill.denom} fillable planted anchors, over ${s.fill.measured} sample(s). The denominator is anchors, not samples.</p>
               ${countTable("anchor left unfilled", s.fill.listed, s.fill.measured)}`
            : "",
          s.syntax
            ? `<div class="cmp-row"><span class="cmp-label">parses</span>${strip(s.syntax.interval)}</div>
               <p class="note">${s.syntax.numer} of ${s.syntax.denom} files checked survived the language's own parser. Syntax only, never correctness.</p>`
            : "",
          s.idempotent
            ? `<div class="cmp-row"><span class="cmp-label">idempotent</span>${strip(s.idempotent.interval)}</div>
               <p class="note">${s.idempotent.numer} of ${s.idempotent.denom} files were unchanged by a second application.</p>`
            : "",
        ].join("") || `<p class="note">No form signal was recorded on this run — it predates them. That is "not recorded", not "nothing parsed".</p>`,
      );

  const expectation = meta?.expect?.actions;
  const selection = s.selection
    ? group(
        "What was selected",
        `<p class="note">A description of the stored answers, not a score. The harness grades selection against this case's expectations at run time and stores the invocations rather than the grade.</p>
        <table class="mini"><thead><tr>
          <th>action</th><th class="num">mean per sample</th><th class="num">samples using it</th>${expectation ? `<th class="num">a complete answer</th>` : ""}
        </tr></thead><tbody>
        ${s.selection.actions
          .map(
            (a) => `<tr><td>${esc(a.action)}</td><td class="num">${a.mean.toFixed(2)}</td><td class="num">${a.samplesSelecting} / ${s.selection.scored}</td>${
              expectation ? `<td class="num">${expectation[a.action] ?? "—"}</td>` : ""
            }</tr>`,
          )
          .join("")}
        </tbody></table>
        ${countTable("chosen first", s.selection.first)}`,
      )
    : "";

  const paths = s.paths
    ? group(
        "Authorized paths written",
        `<div class="cmp-row"><span class="cmp-label">written</span>${strip(s.paths.interval)}</div>
        <p class="note">${s.paths.wrote} of ${s.paths.want} paths the record authorizes. The denominator is paths, not samples.</p>
        ${countTable("never written", s.paths.missing)}
        ${s.paths.unexpected.length ? `<p class="note">Written to a path nothing authorized, and therefore discarded before disk — code at the wrong path is a different finding from no code at all.</p>${countTable("unauthorized path", s.paths.unexpected)}` : ""}`,
      )
    : "";

  const samples = group(
    `Every sample (${s.runs.length})`,
    `<div class="samples">${s.runs
      .map((r, i) => {
        const why = r.outcome === "invalid"
          ? (r.rules || []).join(", ") || "rejected"
          : r.behavior
            ? r.behavior.outcome === "ok"
              ? "application worked"
              : r.behavior.outcome === "checks_failed"
                ? `${r.behavior.passed}/${r.behavior.checks} assertions held`
                : `broke at ${r.behavior.failed_phase || "an unnamed phase"}`
            : "";
        return `<div class="sample">
          <span class="idx">${i + 1}</span>
          <span class="out">${esc(r.outcome)}</span>
          <span class="why">${esc(why)}</span>
          <span class="ms">${fmtDuration(r.ms)}</span>
        </div>`;
      })
      .join("")}</div>
    ${
      s.runs.some((r) => r.invocations?.length)
        ? `<details class="more"><summary>Show the invocations each sample produced</summary>${s.runs
            .map((r, i) =>
              r.invocations?.length
                ? `<p class="note">sample ${i + 1}</p><pre class="inv">${esc(
                    r.invocations.map((v) => `${v.action}(${Object.entries(v.kwargs || {}).map(([k, val]) => `${k}=${JSON.stringify(val)}`).join(", ")})`).join("\n"),
                  )}</pre>`
                : "",
            )
            .join("")}</details>`
        : ""
    }`,
  );

  $("#detail-title").textContent = s.case;
  $("#detail-sub").textContent = [
    fmtDate(s.at),
    s.model.ID,
    s.model.Engine,
    `n=${s.samples}`,
    s.resolution || "unstated",
  ].filter(Boolean).join(" · ");
  $("#detail-inner").innerHTML =
    `${caveats}${outcome}${behaviour}${signals}${selection}${paths}${provenance}${cost}${samples}`;

  const dlg = $("#detail");
  dlg.showModal();
  $("#detail-inner").scrollTop = 0;
}

// ── compare ───────────────────────────────────────────────────────

function renderCompare() {
  const citable = state.rows.filter((r) => r.comparable);
  const label = (r) => `${r.case} · ${r.model.ID} · ${fmtDate(r.at)} · n=${r.samples} · ${r.resolution}`;

  const opts = (sel) =>
    citable
      .map((r, i) => `<option value="${i}"${i === sel ? " selected" : ""}>${esc(label(r))}</option>`)
      .join("");

  if (citable.length < 2) {
    $("#compare-grid").innerHTML = "";
    $("#compare-out").innerHTML = `<p class="note">Fewer than two citable runs are recorded, so there is nothing to compare.</p>`;
    return;
  }

  $("#compare-grid").innerHTML = `
    <span class="ctl"><label for="cmp-a">Run A</label><select id="cmp-a">${opts(0)}</select></span>
    <span class="ctl"><label for="cmp-b">Run B</label><select id="cmp-b">${opts(1)}</select></span>`;

  const draw = () => {
    const a = citable[Number($("#cmp-a").value)];
    const b = citable[Number($("#cmp-b").value)];

    const line = (name, ia, ib, note) => {
      const c = compare(ia, ib);
      return `<div class="dgroup"><h4>${esc(name)}</h4>
        <div class="cmp-row"><span class="cmp-label">A</span>${strip(ia)}</div>
        <div class="cmp-row"><span class="cmp-label">B</span>${strip(ib)}</div>
        <p class="verdict ${c.distinguished ? "" : "overlap"}">${esc(c.text)}</p>
        ${note ? `<p class="note">${esc(note)}</p>` : ""}</div>`;
    };

    // Everything the two runs do not hold in common. Each is a reason the
    // numbers below are not drawn on the same terms, and each is stated rather
    // than left for a reader to notice.
    const mismatches = [];

    if (a.case !== b.case) {
      mismatches.push([
        "different cases",
        "These are two different fixtures — different provenance records, different applications, and possibly " +
          "different generator packages. A difference between them is a difference between two questions, not a " +
          "measurement of one variable.",
      ]);
    }

    if (a.fixture && b.fixture && a.fixture !== b.fixture) {
      mismatches.push([
        "different fixture digests",
        "The packages, records or behaviour target changed between these two runs, so they were not drawn against " +
          "the same inputs.",
      ]);
    } else if (!a.fixture || !b.fixture) {
      mismatches.push([
        "inputs unknown on one side",
        "At least one of these runs predates the fixture digest, so there is no way to confirm the two were drawn " +
          "against the same packages and records. They may or may not have been.",
      ]);
    }

    if (a.retries !== b.retries) {
      mismatches.push([
        "different retry budgets",
        `${a.retries} against ${b.retries}. A raised budget lifts the validity number without lifting the working ` +
          "count, so these validity rates are not drawn on the same terms.",
      ]);
    }

    $("#compare-out").innerHTML =
      `<p class="note"><strong>A</strong> ${esc(label(a))}<br><strong>B</strong> ${esc(label(b))}</p>` +
      mismatches
        .map(([tag, why]) => `<p class="caveat"><strong>${esc(tag)}</strong> — ${esc(why)}</p>`)
        .join("") +
      line("Valid answers", a.validity, b.validity) +
      (a.behavior && b.behavior
        ? line("Applications that work", a.behavior.workingRate, b.behavior.workingRate)
        : `<div class="dgroup"><h4>Applications that work</h4><p class="note">Behaviour was not measured on both runs, so there is nothing to compare.</p></div>`) +
      (a.baseline !== b.baseline
        ? `<p class="note">One of these is the baseline arm — the same model and record with no generator package at all. Selection, binding and anchor fill do not exist on that side, so only the two rows above are comparable.</p>`
        : "");
  };

  $("#cmp-a").addEventListener("change", draw);
  $("#cmp-b").addEventListener("change", draw);
  draw();
}

// ── the worked example in "how to read a rate" ────────────────────

function renderDemo() {
  const rows = [
    [5, 5],
    [10, 10],
    [20, 20],
    [30, 30],
  ];
  $("#demo-intervals").innerHTML =
    rows
      .map(([k, n]) => `<span class="demo-row">${strip(wilson(k, n))}</span>`)
      .join("") +
    `<figcaption>The same perfect score at four sample sizes. Five samples buy far less than they look like they do.</figcaption>`;
}

// ── go ────────────────────────────────────────────────────────────

function render() {
  const rows = visible();
  renderSummary(rows);
  renderTable(rows);
}

load()
  .then(() => {
    renderControls();
    wireSorting();
    renderDemo();
    renderCompare();
    render();
    $("#detail-close").addEventListener("click", () => $("#detail").close());
    $("#detail").addEventListener("click", (e) => {
      if (e.target.id === "detail") e.target.close();
    });
  })
  .catch((err) => {
    $("#runs-body").innerHTML = `<tr><td colspan="9">Could not load the results: ${esc(err)}</td></tr>`;
    console.error(err);
  });
