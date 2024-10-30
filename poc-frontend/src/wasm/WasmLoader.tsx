import './wasm_exec.js';
import './wasmTypes.d.ts';

import React, { useEffect } from 'react';

async function loadWasm(): Promise<void> {
  const goWasm = new window.Go();
  await WebAssembly?.instantiateStreaming(fetch('/assets/doxxier.wasm'), goWasm.importObject).then((result) => {
    goWasm.run(result.instance);
});
}

export const LoadWasm: React.FC<React.PropsWithChildren<{}>> = (props) => {
  const [isLoading, setIsLoading] = React.useState(true);

  useEffect(() => {
    loadWasm().then(() => {
      setIsLoading(false);
    });
  }, []);

  if (isLoading) {
    return (
      <div>
        loading WebAssembly...
      </div>
    );
  } else {
    return <React.Fragment>{props.children}</React.Fragment>;
  }
};