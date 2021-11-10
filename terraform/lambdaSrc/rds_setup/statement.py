createNodeTable = "create table IF NOT EXISTS Nodes (ip varchar(255) NOT NULL, ipString varchar(255), PRIMARY KEY (ip))"
createClusterTable = ("create table IF NOT EXISTS Cluster (idCluster varchar(255) NOT NULL, node varchar(255) NOT NULL, attached int NULL,"
                        " PRIMARY KEY (idCluster, node), INDEX idx0 (node ASC) VISIBLE,"
                        " CONSTRAINT con0 FOREIGN KEY (node) REFERENCES Nodes (ip) ON DELETE CASCADE ON UPDATE NO ACTION)")


dropIfExistsProc = """DROP PROCEDURE IF EXISTS replicaSet"""
replicaProc =     """CREATE PROCEDURE `clusterSet`(in myIP varchar(255),in ipStr varchar(255),in clusterSize int, out valid int, out crashed int)
BEGIN
   
       declare var_nodes int;
   	   declare var_clusterId int;
       declare var_lastCluster int;
       declare var_nodesInCluster int;
       declare var_newClusterId int;

      declare exit handler for sqlexception
      begin
         rollback;
         resignal;
      end;

      set AUTOCOMMIT = 0;
      set session transaction isolation level serializable;
      start transaction;

      transactionBlock:BEGIN

       -- add Node if not exists
      insert ignore into Nodes (ip, ipString) values (myIP, ipStr);

      -- update ipStr
      update Nodes set ipString = ipStr where ip = myIP;

      -- check node number for replication
      select count(ip)
      from Nodes
      into var_nodes;

      -- not enough nodes for first cluster
      if( var_nodes < clusterSize) then
         set valid = 0;
      else
         set valid = 1;
      end if;

      -- check crash if ip is in a cluster
      if(exists(  select node
               from Cluster
               where node = myIP )) then
         set crashed = 1;

         -- get cluster id
         select idCluster
         from Cluster
         where node = myIP into var_clusterId;

   	  -- get ip in the cluster
         select node, ipString
         from Cluster join Nodes on Nodes.ip=Cluster.node
         where node <> myIP and idCluster = var_clusterId;

         leave transactionBlock;
      end if;

      set crashed = 0;

       -- empty list
       if(var_nodes < clusterSize) then
         select ip, ipString from Nodes limit 0;
         leave transactionBlock;
      end if;

      -- first cluster
      if(var_nodes = clusterSize) then
   		insert into Cluster(node, idCluster, attached)
           select ip,0,0 from Nodes;

   		select node, ipString
           from Cluster join Nodes on Nodes.ip=Cluster.node
           where node <> myIP and idCluster = 0;

           leave transactionBlock;
   	end if;

       -- check join or create new cluster
       -- get last cluster
       select max(idCluster)
       from Cluster into var_lastCluster;

       -- get count of nodes attached in last cluster
       select count(*)
       from Cluster
       where idCluster = var_lastCluster and attached=1
       into var_nodesInCluster;

       if(var_nodesInCluster < clusterSize-1) then
   		-- not enough nodes to a new cluster
           -- attach node to the last cluster
           insert into Cluster(idCluster, node, attached)
           values(var_lastCluster, myIP, 1);

           -- return ip of the cluster
           select node, ipString
           from Cluster join Nodes on Nodes.ip=Cluster.node
           where node <> myIP and idCluster = var_lastCluster;
   	else
   		-- create new cluster
           set var_newClusterId = var_lastCluster+1;
           -- insert ip in new cluster
           insert into Cluster(idCluster, node, attached)
           values(var_newClusterId, myIP, 0);

           -- swicth other to my cluster
           update Cluster
           set idCluster = var_newClusterId, attached = 0
           where idCluster = var_lastCluster and attached = 1;

           -- return nodes
           select node, ipString
           from Cluster join Nodes on Nodes.ip=Cluster.node
           where node <> myIP and idCluster = var_newClusterId;

       end if;

      END transactionBlock;
      commit;
   END"""
                     
allOp = [createNodeTable, createClusterTable, dropIfExistsProc, replicaProc]