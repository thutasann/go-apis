import React from 'react';

export default function Post(props) {
  return (
    <div style={{ padding: '2rem', fontFamily: 'system-ui' }}>
      <h1>Post {props.id || 'unknown'}</h1>
      <p>Loaded at: {props.timestamp || 'SSR'}</p>
      <nav style={{ display: 'flex', gap: '1rem' }}>
        <a href='/'>← Home</a>
        <a href='/posts/1'>Post 1</a>
        <a href='/posts/2'>Post 2</a>
      </nav>
    </div>
  );
}

export function getServerSideProps(context) {
  return {
    props: {
      id: context.params ? context.params.id : 'unknown',
      timestamp: new Date().toISOString(),
    },
  };
}
