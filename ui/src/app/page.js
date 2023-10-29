"use client";
import { data } from "autoprefixer";
import { useEffect } from "react";
import { useState } from "react";

function Square({ context }) {
  return (
    <button className="square px-8 border text-stone-50...">{context}</button>
  );
}
var GameText = "";
export default function Home() {
  const [message, SetMessage] = useState([]);
  useEffect(() => {
    fetch("http://localhost:8080/Startgame")
      .then((data) => data.json())
      .then((data) => {
        SetMessage(data);
        GameText = data.message;
        console.log(GameText);
      });
  }, []);

  const [id, setId] = useState("");
  const [data, setData] = useState(null);

  const handleClick = async () => {
    try {
      const data = await (await fetch(`http://localhost:8080/${id}`)).json();
      setData(data);
      console.log(data);
      GameText = data.Announcement;
    } catch (err) {
      console.log(err.message);
      GameText = data.Announcement;
    }

    // try {
    //   const data = await (
    //     await fetch(`http://localhost:8080/Enemy/${id}`)
    //   ).json();
    //   setData(data);
    //   Data = data.name;
    // } catch (err) {
    //   console.log(err.message);
    // }
  };
  return (
    <>
      <div className="left ">
        <div class="relative h-16 ...">
          <div class="absolute inset-x-0 top-2 text-xl ...">
            List of Command
          </div>
        </div>

        <h2 className="ontWhite grid gap-4 grid-cols-2 text-lg">
          <p className="ontWhite"> Fight </p>
          <p className="ontWhite"> SetPlayer/Id </p>
          <p className="ontWhite"> checkStats</p>
          <p className="ontWhite"> GetAreas</p>
          <p className="ontWhite"> </p>
        </h2>
      </div>
      <div className="TopCenter border w-2/4 h-60">{GameText}</div>
      <div className="center ">
        <div className="board-row p-8 flex space-x-40 ...">
          <input
            className="Input"
            required="required"
            placeholder="Enter an ID"
            value={id}
            onChange={(e) => setId(e.target.value)}
          />

          <button type="submit" onClick={handleClick}>
            Search
          </button>
        </div>
      </div>
    </>
  );
}
