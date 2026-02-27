import React from 'react';

export default function Post(props) {
  return (
    <div>
      <h1>Post {props.id || 'unknown'}</h1>
      <p>Dynamic route parameter extracted by Go router.</p>
      <a href='/'>Home</a>
    </div>
  );
}

export function getServerSideProps(context) {
  return {
    props: {
      id: context.params ? context.params.id : 'unknown',
      timestamp: Date.now(),
    },
  };
}
