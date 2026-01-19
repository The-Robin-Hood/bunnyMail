export function TerminalLoader({
  title,
  loadingText,
}: {
  title: string;
  loadingText?: string;
}) {
  return (
    <div className="flex w-full h-screen justify-center items-center bg-black">
      <div className="terminal-loader">
        <div className="terminal-header">
          <div className="terminal-title">{title}</div>
          <div className="terminal-controls">
            <div className="control close"></div>
            <div className="control minimize"></div>
            <div className="control maximize"></div>
          </div>
        </div>
        <div className="text">{loadingText || "Loading..."}</div>
      </div>
    </div>
  );
}

export function RadixLoader() {
  return (
    <div className="banter-loader">
      <div className="banter-loader__box"></div>
      <div className="banter-loader__box"></div>
      <div className="banter-loader__box"></div>
      <div className="banter-loader__box"></div>
      <div className="banter-loader__box"></div>
      <div className="banter-loader__box"></div>
      <div className="banter-loader__box"></div>
      <div className="banter-loader__box"></div>
      <div className="banter-loader__box"></div>
    </div>
  );
}
