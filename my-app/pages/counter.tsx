import React, { useState, useEffect, useCallback } from "react";

const buttonStyle = {
  margin: "5px",
  padding: "10px",
};

type CounterButtonProps = {
  onClick: () => void;
  children: React.ReactNode;
};

function CounterButton({ onClick, children }: CounterButtonProps) {
  return (
    <button style={buttonStyle} onClick={onClick}>
      {children}
    </button>
  );
}

function useCounter(initialValue = 0) {
  const [count, setCount] = useState(initialValue);

  const increment = useCallback(
    () => setCount((prevCount) => prevCount + 1),
    []
  );
  const decrement = useCallback(
    () => setCount((prevCount) => prevCount - 1),
    []
  );

  useEffect(() => {
    // このコードは初期レンダリング後のクライアントサイドでのみ実行されます
    const storedValue = localStorage.getItem("counterValue");
    if (storedValue) {
      setCount(Number(storedValue));
    }
  }, []);

  useEffect(() => {
    localStorage.setItem("counterValue", count.toString());
  }, [count]);

  return { count, increment, decrement };
}

function Counter() {
  // デフォルトの値で初期化します；存在する場合は保存された値で更新されます
  const { count, increment, decrement } = useCounter(0);

  return (
    <div>
      <p>カウント: {count}</p>
      <CounterButton onClick={increment}>増加</CounterButton>
      <CounterButton onClick={decrement}>減少</CounterButton>
    </div>
  );
}

export default Counter;
