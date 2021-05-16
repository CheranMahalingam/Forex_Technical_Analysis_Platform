import React from "react";
import { useHistory } from "react-router-dom";
import { makeStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import "./HomePage.css";
import NavigationMenu from "../components/NavigationMenu/NavigationMenu";

const useStyles = makeStyles({
  buttonGroup: {
    marginLeft: "36%",
    marginTop: "10%",
  },
  button1: {
    position: "relative",
    background: "linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)",
    border: 0,
    borderRadius: 3,
    boxShadow: "0 3px 5px 2px rgba(33, 203, 243, .3)",
    color: "white",
    height: 48,
    padding: "0 30px",
    width: 190,
  },
  button2: {
    position: "absolute",
    background: "linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)",
    border: 0,
    borderRadius: 3,
    boxShadow: "0 3px 5px 2px rgba(33, 203, 243, .3)",
    color: "white",
    height: 48,
    padding: "0 30px",
    width: 190,
    marginLeft: 50
  },
});

function HomePage() {
  const classes = useStyles();
  const history = useHistory();

  return (
    <div className="background">
      <NavigationMenu />
      <h1 className="title">
        Forex Analytics
      </h1>
      <div className={classes.buttonGroup}>
        <Button
          className={classes.button1}
          onClick={() => history.push("/charts")}
        >
          Exchange Rates
        </Button>
        <Button
          className={classes.button2}
          onClick={() => history.push("/news")}
        >
          Market News
        </Button>
      </div>
    </div>
  );
}

export default HomePage;
