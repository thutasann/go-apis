import React from 'react';

export default function Home(props) {
  return (
    <div style={{ padding: '2rem', fontFamily: 'system-ui' }}>
      <h1>reactgo</h1>
      <p>Server-rendered by Go + V8. SPA navigation enabled.</p>
      <nav style={{ display: 'flex', gap: '1rem' }}>
        <a href='/about'>About</a>
        <a href='/posts/42'>Post 42</a>
        <a href='/posts/99'>Post 99</a>
      </nav>
      <p style={{ color: '#666', marginTop: '1rem' }}>Click links above — no full page reload. Check Network tab.</p>
    </div>
  );
}
