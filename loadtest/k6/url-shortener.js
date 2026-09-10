import http from "k6/http";
import { check } from "k6";

const baseURL = __ENV.BASE_URL || "http://localhost:8080";

export const options = {
  scenarios: {
    safe_steady_load: {
      executor: "constant-arrival-rate",
      rate: Number(__ENV.RATE || 2),
      timeUnit: "1s",
      duration: __ENV.DURATION || "1m",
      preAllocatedVUs: Number(__ENV.PRE_ALLOCATED_VUS || 4),
      maxVUs: Number(__ENV.MAX_VUS || 8),
    },
  },
  thresholds: {
    checks: ["rate>0.99"],
    http_req_failed: ["rate<0.01"],
    "http_req_duration{endpoint:create_short_url}": ["p(95)<500"],
    "http_req_duration{endpoint:redirect_short_url}": ["p(95)<500"],
  },
};

function targetURL() {
  const suffix = `${__VU}-${__ITER}-${Date.now()}-${Math.random().toString(16).slice(2)}`;
  return `https://example.com/load-test/${suffix}`;
}

function extractCode(shortURL) {
  if (!shortURL) {
    return "";
  }

  const parts = shortURL.split("/").filter(Boolean);
  return parts[parts.length - 1] || "";
}

export default function () {
  const createRes = http.post(
    `${baseURL}/urls/`,
    JSON.stringify({ url: targetURL() }),
    {
      headers: {
        "Content-Type": "application/json",
      },
      tags: {
        endpoint: "create_short_url",
        route: "/urls",
      },
    },
  );

  const created = check(createRes, {
    "create returns 201": (res) => res.status === 201,
    "create returns short url": (res) => Boolean(res.json("url")),
  });

  if (!created) {
    return;
  }

  const code = extractCode(createRes.json("url"));
  const hasCode = check(code, {
    "short code is present": (value) => value.length > 0,
  });

  if (!hasCode) {
    return;
  }

  const redirectRes = http.get(`${baseURL}/${code}`, {
    redirects: 0,
    tags: {
      endpoint: "redirect_short_url",
      route: "/{code}",
    },
  });

  check(redirectRes, {
    "redirect returns 302": (res) => res.status === 302,
    "redirect location is present": (res) => Boolean(res.headers.Location),
  });
}
