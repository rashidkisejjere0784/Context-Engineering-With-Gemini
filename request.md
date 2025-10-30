# Feature Request

## 🎯 The Goal

Build a **simple web-based tool** (single-page app) in HTML/CSS/JavaScript that allows a user to get the current weather for a specified city.

The web app should:

1. Provide a friendly, modern UI (search input + results card).
2. Use the free Open-Meteo Geocoding API to resolve a city name to latitude/longitude.
3. Use the Open-Meteo Weather API to fetch the `current_weather` for the coordinates.
4. Display the temperature and a few useful fields (wind speed, weather code description, timestamp) in a clean card.
5. Handle errors gracefully (city not found, network / API errors, rate limits) and show helpful messages in the UI.
6. Be easily testable locally (single HTML file that works without a backend) and production-ready with optional static hosting instructions (GitHub Pages / Netlify).

> No API keys required — Open-Meteo is free and public for these endpoints.

---

## 🎨 UX / UI Guidance

Keep it minimal and accessible — mobile-first design.

**Primary screen elements:**

* Header with app title and small description.
* Search box with placeholder `"Enter city (e.g. Kampala)"` and a "Search" button.
* Recent searches list (localStorage-backed) to speed repeat lookups.
* Weather result card showing:

  * City name + country
  * Current temperature (°C) and a friendly icon/emoji for the weather code
  * Wind speed and direction
  * Time of the observation (localized)
  * A small note describing the data source (Open-Meteo)
* Inline error / status area for messages ("Searching...", "City not found", etc.)

**Visual style:** clean card UI, centered content, subtle shadows, 2xl rounded corners. Small animations for transitions.

Accessibility: keyboard-first, labels for the input, ARIA live region for dynamic status messages, good color contrast.

---

## 📚 Required APIs & Endpoints

* **Geocoding API** (Open-Meteo) — find latitude & longitude by city name.

  * Example: `https://geocoding-api.open-meteo.com/v1/search?name=Kampala&count=1`

* **Weather API** (Open-Meteo) — get current weather using `latitude` and `longitude`.

  * Example: `https://api.open-meteo.com/v1/forecast?latitude=0.3166&longitude=32.5825&current_weather=true&timezone=auto`

We only need `current_weather.temperature`, `current_weather.windspeed`, `current_weather.winddirection`, and `current_weather.time`.

---

## 🧩 Technical Requirements

* Single-page app: `index.html` with embedded CSS and JS (for easy testing).
* Plain vanilla JavaScript (no frameworks required). Keep code small and readable.
* Use `fetch()` for requests and `async/await`.
* Save the last 5 searches in `localStorage` and allow click-to-search.
* Show a loading state while fetching.
* Unit-testable JavaScript functions (e.g. `getCoordinates(city)`, `getCurrentWeather(lat, lon)`, `formatWeatherData()`).

---

## ⚠️ Potential Pitfalls & Gotchas

* Geocoding may return multiple matches — the app should select the top result by default and allow the user to pick another if provided.
* Network errors and CORS: Open-Meteo supports CORS so front-end fetches should work from static hosting.
* Timezones: use `timezone=auto` to let the API return times localized to the requested coordinates.
* Rate limits: keep requests reasonable and show an explanatory message if the API returns an error status.

---

## ✅ Acceptance Criteria

1. Visiting `index.html` shows a search input and a help text.
2. Typing a valid city and pressing Search shows the current temperature and other fields.
3. Errors display human-friendly messages (e.g., "City not found", "Network error — try again").
4. Recent searches persist across reloads and are clickable.
5. No API keys are necessary.
6. Code is commented and modular enough to unit-test the fetch functions.

---

## 📦 Local Development & Testing

* Chrome / Firefox: open `index.html` directly (file URI) or use a tiny static server (recommended) like `python -m http.server`.
* Tests: Provide simple unit tests with your favorite runner (Jest for node-based testing or plain assertions in a test HTML page). Tests should stub fetch responses.

---

## 🔧 Example: `index.html` Prototype (single-file)

Below is a compact, production-adjacent single-file prototype you can drop into a project and open in a browser.

```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <title>Web Weather — ESM Prototype</title>
  <style>
    :root{font-family:Inter,system-ui,Segoe UI,Roboto,'Helvetica Neue',Arial;}
    body{display:flex;min-height:100vh;align-items:center;justify-content:center;background:#f3f4f6;margin:0;padding:24px}
    .container{width:100%;max-width:720px;background:#fff;border-radius:16px;box-shadow:0 8px 30px rgba(2,6,23,0.08);padding:24px}
    header{display:flex;align-items:center;justify-content:space-between}
    h1{margin:0;font-size:20px}
    .search{display:flex;gap:8px;margin-top:18px}
    input[type=text]{flex:1;padding:12px;border-radius:10px;border:1px solid #e6e9ef;font-size:16px}
    button{padding:10px 14px;border-radius:10px;border:0;background:#0ea5a3;color:#fff;font-weight:600}
    .card{margin-top:18px;padding:16px;border-radius:12px;background:linear-gradient(180deg,#ffffff,#fbfdff);box-shadow:0 2px 8px rgba(2,6,23,0.04)}
    .small{font-size:13px;color:#6b7280}
    .error{color:#dc2626}
    .recent{margin-top:12px;display:flex;gap:8px;flex-wrap:wrap}
    .chip{background:#eef2ff;padding:6px 10px;border-radius:999px;cursor:pointer}
  </style>
</head>
<body>
  <div class="container" role="main">
    <header>
      <h1>Web Weather</h1>
      <div class="small">Powered by Open-Meteo — no API key needed</div>
    </header>

    <label for="city" class="small" style="margin-top:12px;display:block">Enter a city name</label>
    <div class="search">
      <input id="city" type="text" placeholder="e.g. Kampala" aria-label="City name" />
      <button id="searchBtn">Search</button>
    </div>

    <div id="status" class="small" role="status" aria-live="polite" style="margin-top:10px"></div>

    <div id="recent" class="recent" aria-label="Recent searches"></div>

    <div id="result" class="card" style="display:none"></div>
  </div>

  <script>
    // Helper: map Open-Meteo weather codes to emoji descriptions (small friendly shortcut)
    const weatherCodeMap = {
      0: ['Clear', '☀️'],
      1: ['Mainly clear', '🌤️'],
      2: ['Partly cloudy', '⛅'],
      3: ['Overcast', '☁️'],
      45: ['Fog', '🌫️'],
      48: ['Depositing rime fog', '🌫️'],
      51: ['Light drizzle', '🌧️'],
      53: ['Moderate drizzle', '🌧️'],
      55: ['Dense drizzle', '🌧️'],
      61: ['Slight rain', '🌧️'],
      63: ['Moderate rain', '🌧️'],
      80: ['Rain showers', '🌦️'],
      95: ['Thunderstorm', '⛈️'],
    };

    // DOM
    const cityInput = document.getElementById('city');
    const searchBtn = document.getElementById('searchBtn');
    const statusEl = document.getElementById('status');
    const resultEl = document.getElementById('result');
    const recentEl = document.getElementById('recent');

    // Local storage keys
    const RECENT_KEY = 'web-weather-recent';

    // Utils
    function setStatus(text, isError=false){ statusEl.textContent = text; statusEl.className = isError ? 'small error' : 'small'; }

    async function getCoordinates(city){
      const url = `https://geocoding-api.open-meteo.com/v1/search?name=${encodeURIComponent(city)}&count=1`;
      const res = await fetch(url);
      if(!res.ok) throw new Error('Geocoding API error');
      const json = await res.json();
      if(!json.results || json.results.length === 0) return null;
      return json.results[0];
    }

    async function getCurrentWeather(lat, lon){
      const url = `https://api.open-meteo.com/v1/forecast?latitude=${lat}&longitude=${lon}&current_weather=true&timezone=auto`;
      const res = await fetch(url);
      if(!res.ok) throw new Error('Weather API error');
      return res.json();
    }

    function renderResult(place, data){
      const cw = data.current_weather;
      const [desc, emoji] = weatherCodeMap[cw.weathercode] || ['Unknown','🌈'];
      const time = new Date(cw.time).toLocaleString();
      resultEl.style.display = 'block';
      resultEl.innerHTML = `
        <div style="display:flex;justify-content:space-between;align-items:center">
          <div>
            <div style="font-weight:700;font-size:18px">${place.name}, ${place.country}</div>
            <div class="small">Observation: ${time}</div>
          </div>
          <div style="text-align:right">
            <div style="font-size:32px;font-weight:800">${cw.temperature}°C</div>
            <div class="small">${emoji} ${desc}</div>
          </div>
        </div>
        <hr style="margin:12px 0">
        <div class="small">Wind: ${cw.windspeed} m/s at ${cw.winddirection}°</div>
        <div class="small" style="margin-top:8px">Data from Open-Meteo</div>
      `;
    }

    function saveRecent(city){
      try{
        const list = JSON.parse(localStorage.getItem(RECENT_KEY) || '[]');
        const trimmed = city.trim();
        if(!trimmed) return;
        const dedup = [trimmed, ...list.filter(s => s.toLowerCase() !== trimmed.toLowerCase())];
        localStorage.setItem(RECENT_KEY, JSON.stringify(dedup.slice(0,5)));
        renderRecent();
      }catch(e){console.warn(e)}
    }

    function renderRecent(){
      const list = JSON.parse(localStorage.getItem(RECENT_KEY) || '[]');
      recentEl.innerHTML = '';
      list.forEach(s => {
        const btn = document.createElement('button');
        btn.className = 'chip';
        btn.type = 'button';
        btn.textContent = s;
        btn.addEventListener('click', () => { cityInput.value = s; doSearch(s); });
        recentEl.appendChild(btn);
      })
    }

    async function doSearch(overrideCity){
      const city = (overrideCity || cityInput.value).trim();
      if(!city) { setStatus('Please enter a city name', true); return; }
      setStatus('Searching...');
      resultEl.style.display = 'none';
      try{
        const place = await getCoordinates(city);
        if(!place){ setStatus('City not found — try a different name', true); return; }
        setStatus('Loading weather...');
        const weather = await getCurrentWeather(place.latitude, place.longitude);
        renderResult(place, weather);
        saveRecent(city);
        setStatus('');
      }catch(err){
        console.error(err);
        setStatus('Network or API error — please try again', true);
      }
    }

    // Bind events
    searchBtn.addEventListener('click', () => doSearch());
    cityInput.addEventListener('keydown', (e) => { if(e.key === 'Enter') doSearch(); });

    // Initialize
    renderRecent();
  </script>
</body>
</html>
```

---

## 🧪 Minimal Test Cases (manual)

1. Search `Kampala` — should return a result and show temperature.
2. Search typo / nonsense like `asdasdasd` — should show "City not found".
3. Turn off network -> search -> verify network error message.
4. Click recent item -> check result updates.

---

## 🚀 Deploying

* For quick sharing: push `index.html` to GitHub and enable GitHub Pages (or drag the file into Netlify Drop).
* The app is static and requires no server-side code.

---

## 🧭 Next Improvements (optional)

* Add a small unit-test harness to test `getCoordinates` and `getCurrentWeather` using stubs/mocks.
* Add a map preview (static map) for the returned coordinates.
* Allow user to pick alternative geocoding matches when multiple results are returned.
* Add localization (°F toggle, language strings).
