import React from 'react';

export default function Home(props) {
  return (
    <div>
      <h1>reactgo</h1>
      <p>Server-rendered by Go + V8</p>
      <a href='/about'>About</a>
      <a href='/posts/42'>Post 42</a>
    </div>
  );
}
