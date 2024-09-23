import {newGame as newGame, createCard} from "./bjgame.js";
window.onload = () => {
  // create a couple of elements in an otherwise empty HTML page
  // const heading = document.createElement("h1");
  // const headingText = document.createTextNode("Big Head!");
  createCard("test","meow1");
  createCard("test","meow2");
  document.querySelector("#makecard").addEventListener('click',createCard);
  document.querySelector("#newgame").addEventListener('click',newGame);

};
