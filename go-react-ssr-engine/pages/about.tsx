import React from 'react';
import Button from '../components/button';

export default function About(props) {
  return (
    <div>
      <h1>About Page</h1>
      <p>A Go-powered React SSR engine.</p>
      <a href='/'>Home</a>
      <Button />
    </div>
  );
}
