import React, { useState, useEffect } from "react";
import { parseNewMarketNews } from "../../utils/websocketMessage";
import NewsCard from "./NewsCard";

var socket;

function MarketNews(props) {
  const [marketNews, setMarketNews] = useState({
    timestamp: [],
    headline: [],
    image: [],
    source: [],
    summary: [],
    url: [],
  });

  const createNews = (news) => {
    const parsedNews = JSON.parse(news);
    const newNews = parseNewMarketNews(parsedNews["News"], { ...marketNews });
    setMarketNews(newNews);
  };

  const connectToNewsSocket = () => {
    socket.send(JSON.stringify({ message: "subscribe", data: "MarketNews" }));
  };

  useEffect(() => {
    socket = new WebSocket(process.env.REACT_APP_WEBSOCKET_URI);

    socket.onopen = () => {
      connectToNewsSocket();
      props.handleEndLoading();
      console.log("websocket connected");
    };

    socket.onclose = () => {
      console.log("websocket disconnected");
    };

    socket.onmessage = (evt) => {
      console.log("onmessage");
      createNews(evt.data);
    };

    socket.onerror = (err) => {
      console.log("encountered error:", err);
      socket.close();
    };

    return () => {
      console.log("component unmounted");
      socket.close();
    };
  }, []);

  const createNewsCard = (column) => {
    let columnNews = {};
    columnNews.timestamp = marketNews.timestamp.filter(
      (_, index) => index % 2 === column
    );
    columnNews.headline = marketNews.headline.filter(
      (_, index) => index % 2 === column
    );
    columnNews.image = marketNews.image.filter(
      (_, index) => index % 2 === column
    );
    columnNews.summary = marketNews.summary.filter(
      (_, index) => index % 2 === column
    );
    columnNews.source = marketNews.source.filter(
      (_, index) => index % 2 === column
    );
    columnNews.url = marketNews.url.filter((_, index) => index % 2 === column);

    return columnNews.timestamp.map((_, index) => (
      <NewsCard
        key={columnNews.timestamp[index]}
        timestamp={columnNews.timestamp[index]}
        title={columnNews.headline[index]}
        image={columnNews.image[index]}
        summary={columnNews.summary[index]}
        source={columnNews.source[index]}
        newsurl={columnNews.url[index]}
      />
    ));
  };

  return (
    <React.Fragment>
      <div style={{ display: "flex" }}>
        <div style={{ flex: "50%" }}>
          {createNewsCard(0)}
        </div>
        <div style={{ flex: "50%" }}>
          {createNewsCard(1)}
        </div>
      </div>
    </React.Fragment>
  );
}

export default MarketNews;
