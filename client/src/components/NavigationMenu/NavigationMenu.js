import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import Drawer from "@material-ui/core/Drawer";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import Button from "@material-ui/core/Button";
import MenuIcon from "@material-ui/icons/Menu";
import HomeIcon from "@material-ui/icons/Home";
import TrendingUpIcon from "@material-ui/icons/TrendingUp";
import PublicIcon from "@material-ui/icons/Public";
import EventIcon from "@material-ui/icons/Event";
import InfoIcon from "@material-ui/icons/Info";
import GitHubIcon from "@material-ui/icons/GitHub";

function NavigationMenu() {
  const [openMenu, setOpenMenu] = useState(false);

  const history = useHistory();

  const handleMenuClose = () => {
    setOpenMenu(false);
  };

  const handleMenuToggle = () => {
    setOpenMenu(!openMenu);
  };

  const menuItemList = () => {
    return (
      <div onClick={handleMenuClose} style={{ width: 300 }}>
        <List>
          <ListItem
            button
            onClick={() => history.push("/")}
            style={{ marginTop: "12%" }}
          >
            <ListItemIcon>
              <HomeIcon />
            </ListItemIcon>
            <ListItemText primary="Home" />
          </ListItem>
          <ListItem
            button
            onClick={() => history.push("/charts")}
            style={{ marginTop: "10%" }}
          >
            <ListItemIcon>
              <TrendingUpIcon />
            </ListItemIcon>
            <ListItemText primary="Charts" />
          </ListItem>
          <ListItem
            button
            onClick={() => history.push("/news")}
            style={{ marginTop: "10%" }}
          >
            <ListItemIcon>
              <PublicIcon />
            </ListItemIcon>
            <ListItemText primary="Market News" />
          </ListItem>
          {/* <ListItem
            button
            onClick={() => history.push("/calendar")}
            style={{ marginTop: "10%" }}
          >
            <ListItemIcon>
              <EventIcon />
            </ListItemIcon>
            <ListItemText primary="Economic Calendar" />
          </ListItem> */}
          {/* <ListItem button style={{ marginTop: "10%" }}>
            <ListItemIcon>
              <InfoIcon />
            </ListItemIcon>
            <ListItemText primary="About" />
          </ListItem> */}
          <a
            href="https://github.com/CheranMahalingam/Forex_Technical_Analysis_Platform"
            style={{ color: "inherit", textDecoration: "none" }}
          >
            <ListItem button style={{ marginTop: "12%" }}>
              <ListItemIcon>
                <GitHubIcon />
              </ListItemIcon>
              <ListItemText primary="Source Code" />
            </ListItem>
          </a>
        </List>
      </div>
    );
  };

  return (
    <React.Fragment>
      <Button
        onClick={handleMenuToggle}
        startIcon={<MenuIcon />}
        size="large"
        style={{
          position: "relative",
          borderRadius: 50,
          marginLeft: 20,
          marginTop: 20,
          color: "#FFFFFF",
        }}
      />
      <Drawer anchor="left" open={openMenu} onClose={handleMenuClose}>
        {menuItemList()}
      </Drawer>
    </React.Fragment>
  );
}

export default NavigationMenu;
