import Image from "next/image";
import Container from "postcss/lib/container";

function Square({ context }) {
  return (
    <button className="square px-8 border text-stone-50...">{context}</button>
  );
}
export default function Home() {
  return (
    <>
      <div className="center ">
        <div className="board-row p-8 flex space-x-40 ...">
          <Square class="AttackB" context="Attack" />
          <Square class="ItemB" context="Item" />
          <Square context="Check" />
        </div>
      </div>
      <div class="right flex flex-col space-y-4 ...">
        <Square context="Check" />
        <Square context="Check" />
        <Square context="Check" />
      </div>
    </>
  );
}
