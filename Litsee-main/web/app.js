const state = {
  token: localStorage.getItem('litsee_demo_token') || '',
};

const el = {
  notifications: document.getElementById('notifications'),
  registerForm: document.getElementById('register-form'),
  loginForm: document.getElementById('login-form'),
  profileForm: document.getElementById('profile-form'),
  tokenPreview: document.getElementById('token-preview'),
  log: document.getElementById('log'),
  healthStatus: document.getElementById('health-status'),
  profileJson: document.getElementById('profile-json'),
  booksList: document.getElementById('books-list'),
  activityList: document.getElementById('activity-list'),
  booksLimit: document.getElementById('books-limit'),
  booksOffset: document.getElementById('books-offset'),
  activityLimit: document.getElementById('activity-limit'),
  activityOffset: document.getElementById('activity-offset'),
};

function logLine(message, ok = true) {
  const div = document.createElement('div');
  div.className = `log-line ${ok ? 'log-ok' : 'log-bad'}`;
  div.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
  el.log.append(div);
  el.log.scrollTop = el.log.scrollHeight;
}

function notify(message, ok = true) {
  const toast = document.createElement('div');
  toast.className = `toast ${ok ? 'toast-ok' : 'toast-bad'}`;
  toast.textContent = message;
  el.notifications.prepend(toast);

  window.setTimeout(() => {
    toast.remove();
  }, 2600);
}

function setToken(token) {
  state.token = token || '';
  if (state.token) {
    localStorage.setItem('litsee_demo_token', state.token);
    el.tokenPreview.textContent = `${state.token.slice(0, 24)}...`;
  } else {
    localStorage.removeItem('litsee_demo_token');
    el.tokenPreview.textContent = 'не авторизован';
  }
}

async function api(path, { method = 'GET', body, auth = false } = {}) {
  const headers = {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  };

  if (auth) {
    if (!state.token) {
      throw new Error('Нужен JWT. Сначала выполните вход.');
    }
    headers.Authorization = `Bearer ${state.token}`;
  }

  const response = await fetch(path, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  const raw = await response.text();
  let payload = raw;
  if (raw) {
    try {
      payload = JSON.parse(raw);
    } catch (_) {
      payload = raw;
    }
  }

  if (!response.ok) {
    const message = typeof payload === 'string'
      ? payload
      : payload?.error || payload?.message || JSON.stringify(payload);
    throw new Error(`${response.status}: ${message}`);
  }

  return payload;
}

function fillProfileForm(profile) {
  el.profileForm.first_name.value = profile.first_name || '';
  el.profileForm.last_name.value = profile.last_name || '';
  el.profileForm.avatar.value = profile.avatar || '';
  el.profileForm.bio.value = profile.bio || '';
  el.profileJson.textContent = JSON.stringify(profile, null, 2);
}

function renderBooks(books) {
  el.booksList.innerHTML = '';
  if (!books.length) {
    el.booksList.innerHTML = '<div class="item">Пока нет купленных книг.</div>';
    return;
  }

  for (const book of books) {
    const item = document.createElement('div');
    item.className = 'item';
    item.innerHTML = `
      <div class="item-title">${book.title || `Book #${book.book_id}`}</div>
      <div class="item-sub">book_id: ${book.book_id}</div>
      <div class="item-sub">author: ${book.author || '-'}</div>
      <div class="item-sub">order: ${book.order_id || '-'}</div>
      <div class="item-sub">purchased_at: ${book.purchased_at || '-'}</div>
    `;
    el.booksList.append(item);
  }
}

function renderActivity(items) {
  el.activityList.innerHTML = '';
  if (!items.length) {
    el.activityList.innerHTML = '<div class="item">Нет записей активности.</div>';
    return;
  }

  for (const row of items) {
    const item = document.createElement('div');
    item.className = 'item';
    item.innerHTML = `
      <div class="item-title">${row.action}</div>
      <div class="item-sub">status: ${row.status || '-'}</div>
      <div class="item-sub">details: ${row.details || '-'}</div>
      <div class="item-sub">created_at: ${row.created_at || '-'}</div>
    `;
    el.activityList.append(item);
  }
}

async function checkHealth() {
  const checks = [];

  try {
    await api('/api/auth/login', {
      method: 'POST',
      body: { email: 'healthcheck@example.com', password: 'invalid' },
      auth: false,
    });
    checks.push('auth: ok');
  } catch (err) {
    if (String(err.message).startsWith('401:')) {
      checks.push('auth: ok (reachable)');
    } else {
      checks.push(`auth: fail (${err.message})`);
    }
  }

  try {
    await api('/api/profile', { method: 'GET', auth: false });
    checks.push('profile: ok');
  } catch (err) {
    if (String(err.message).startsWith('401:')) {
      checks.push('profile: ok (reachable)');
    } else {
      checks.push(`profile: fail (${err.message})`);
    }
  }

  const allOk = checks.every((line) => line.includes('ok'));
  el.healthStatus.textContent = allOk ? 'сервисы доступны' : 'есть проблемы';
  logLine(`Health check -> ${checks.join(' | ')}`, allOk);
}

el.registerForm.addEventListener('submit', async (event) => {
  event.preventDefault();
  const form = new FormData(el.registerForm);

  try {
    const payload = {
      email: String(form.get('email') || '').trim(),
      password: String(form.get('password') || ''),
    };

    await api('/api/auth/register', { method: 'POST', body: payload });
    logLine(`Регистрация успешна: ${payload.email}`);
    notify(`Регистрация: ${payload.email}`, true);
  } catch (err) {
    logLine(`Регистрация не удалась: ${err.message}`, false);
    notify(`Ошибка регистрации: ${err.message}`, false);
  }
});

el.loginForm.addEventListener('submit', async (event) => {
  event.preventDefault();
  const form = new FormData(el.loginForm);

  try {
    const payload = {
      email: String(form.get('email') || '').trim(),
      password: String(form.get('password') || ''),
    };

    const data = await api('/api/auth/login', { method: 'POST', body: payload });
    if (!data || typeof data !== 'object' || !data.token) {
      throw new Error('JWT не получен от auth-сервиса');
    }
    setToken(data.token);
    logLine(`Вход выполнен: ${payload.email}`);
    notify(`Вход выполнен: ${payload.email}`, true);
  } catch (err) {
    logLine(`Ошибка входа: ${err.message}`, false);
    notify(`Ошибка входа: ${err.message}`, false);
  }
});

document.getElementById('logout').addEventListener('click', () => {
  setToken('');
  logLine('JWT очищен, выполнен выход');
  notify('Выход выполнен', true);
});

document.getElementById('check-health').addEventListener('click', () => {
  checkHealth();
});

document.getElementById('load-profile').addEventListener('click', async () => {
  try {
    const profile = await api('/api/profile', { auth: true });
    fillProfileForm(profile);
    logLine('Профиль загружен');
    notify('Профиль загружен', true);
  } catch (err) {
    logLine(`Не удалось загрузить профиль: ${err.message}`, false);
    notify(`Ошибка профиля: ${err.message}`, false);
  }
});

document.getElementById('save-profile').addEventListener('click', async () => {
  try {
    const payload = {
      first_name: el.profileForm.first_name.value.trim(),
      last_name: el.profileForm.last_name.value.trim(),
      avatar: el.profileForm.avatar.value.trim(),
      bio: el.profileForm.bio.value.trim(),
    };

    await api('/api/profile', { method: 'PUT', body: payload, auth: true });
    logLine('Профиль сохранен');
    notify('Профиль сохранен', true);

    const updated = await api('/api/profile', { auth: true });
    fillProfileForm(updated);
  } catch (err) {
    logLine(`Ошибка сохранения профиля: ${err.message}`, false);
    notify(`Ошибка сохранения: ${err.message}`, false);
  }
});

document.getElementById('load-books').addEventListener('click', async () => {
  try {
    const limit = Number(el.booksLimit.value || 50);
    const offset = Number(el.booksOffset.value || 0);
    const books = await api(`/api/profile/books?limit=${limit}&offset=${offset}`, { auth: true });
    renderBooks(Array.isArray(books) ? books : []);
    logLine(`Купленные книги загружены: ${Array.isArray(books) ? books.length : 0}`);
    notify(`Книги загружены: ${Array.isArray(books) ? books.length : 0}`, true);
  } catch (err) {
    logLine(`Ошибка загрузки книг: ${err.message}`, false);
    notify(`Ошибка книг: ${err.message}`, false);
  }
});

document.getElementById('load-activity').addEventListener('click', async () => {
  try {
    const limit = Number(el.activityLimit.value || 20);
    const offset = Number(el.activityOffset.value || 0);
    const history = await api(`/api/profile/activity?limit=${limit}&offset=${offset}`, { auth: true });
    renderActivity(Array.isArray(history) ? history : []);
    logLine(`История активности загружена: ${Array.isArray(history) ? history.length : 0}`);
    notify(`Активность загружена: ${Array.isArray(history) ? history.length : 0}`, true);
  } catch (err) {
    logLine(`Ошибка загрузки активности: ${err.message}`, false);
    notify(`Ошибка активности: ${err.message}`, false);
  }
});

setToken(state.token);
logLine('Web demo готов. Проверьте сервисы и выполните login.');
