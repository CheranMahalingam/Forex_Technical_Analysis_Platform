import React, { useState } from "react";
import Backdrop from "@material-ui/core/Backdrop";
import CircularProgress from "@material-ui/core/CircularProgress";
import { makeStyles } from "@material-ui/core/styles";
import MarketNews from "../components/MarketNews/MarketNews";
import NavigationMenu from "../components/NavigationMenu/NavigationMenu";

const useStyles = makeStyles((theme) => ({
  backdrop: {
    zIndex: theme.zIndex.drawer + 1,
    color: "#fff",
  },
}));

function NewsPage() {
  const [loading, setLoading] = useState(true);

  const classes = useStyles();

  const handleWebsocketOpen = () => {
    setLoading(false);
  };

  return (
    <React.Fragment>
      <Backdrop className={classes.backdrop} open={loading}>
        <CircularProgress color="inherit" />
      </Backdrop>
      <NavigationMenu />
      <MarketNews handleEndLoading={handleWebsocketOpen} />
    </React.Fragment>
  );
}

export default NewsPage;
