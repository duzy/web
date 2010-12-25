<!DOCTYPE html>
<html>
  <head>
    <title>Dusell.com CPanel</title>
    <link rel="stylesheet" href="/css/cpanel.css" />
  </head>
  <body>
    <div id="page">
      <div id="header">
        <a href="{script}">Home</a> | <a href="{base}/view">view</a>
        <h1>CPanel</h1>
      </div>
      <div id="content-left">
        <ul>
          <li><a href="{script}/eBay">eBay</a></li>
          <hr/>
          <li><a href="{script}/sync">Sync</a></li>
          <hr/>
          <li><a href="{script}/categories">Categories</a></li>
          <li><a href="{script}/items">Items</a></li>
          <li><a href="{script}/orders">Orders</a></li>
          <li><a href="{script}/sells">Selling</a></li>
        </ul>
      </div>
      <div id="content">{content}</div>
      <div id="content-right">{rightside}</div>
      <div id="footer">footer</div>
    </div>
    <div id="page-tail"/>
  </body>
</html>
