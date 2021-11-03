createNodeTable = "create table IF NOT EXISTS Nodes (nodeIP varchar(255) NOT NULL, ipString varchar(255), replicaCount int, PRIMARY KEY (nodeIP))"
createReplicaTable = ("create table IF NOT EXISTS ReplicaOf (Master varchar(255) NOT NULL, Replica varchar(255) NOT NULL,"
                        " PRIMARY KEY (Master, Replica), INDEX idx0 (Replica ASC) VISIBLE, INDEX idx1 (Master ASC) VISIBLE,"
                        " CONSTRAINT con0 FOREIGN KEY (Master) REFERENCES Nodes (nodeIP) ON DELETE CASCADE ON UPDATE NO ACTION,"
                        " CONSTRAINT con1 FOREIGN KEY (Replica) REFERENCES Nodes (nodeIP) ON DELETE CASCADE ON UPDATE NO ACTION)")


dropIfExistsRep = """DROP PROCEDURE IF EXISTS replicaSet"""
replicaProc =     """CREATE PROCEDURE `replicaSet`(in myIP varchar(255),in ipStr varchar(255),in replicaNumber int, out valid int, out crashed int)
BEGIN
   
   declare var_nodes int;
   declare n int default 0;
   declare i int default 0;
    
   declare exit handler for sqlexception
   begin
      rollback;
      resignal;
   end;
   
   drop temporary table if exists tempNode;
   create temporary table tempNode(repIP varchar(255), repIpStr varchar(255));
   
   set AUTOCOMMIT = 0; 
   set session transaction isolation level serializable;
   start transaction;                        
   
   transactionBlock:BEGIN
   
    -- add Node if not exists
   insert ignore into Nodes (nodeIP, ipString, replicaCount) values (myIP, ipStr, 0);

   -- update ipStr
   update Nodes set ipString = ipStr where nodeIP = myIP;

   -- check node number for replication
   select count(nodeIP)
   from Nodes
   where nodeIP <> myIP 
   into var_nodes;

   -- not enough nodes 
   if( var_nodes < replicaNumber) then
      set valid = 0;
   else 
      set valid = 1;
   end if;

   -- check crash
   if(exists(  select Master
            from ReplicaOf
            where Master = myIP )) then
      set crashed = 1;

      -- return old replicas
      select Replica, ipString
      from ReplicaOf join Nodes on Nodes.nodeIP=ReplicaOf.Replica
      where Master = myIP; 
                           
      leave transactionBlock;
   end if;
                        
   set crashed = 0;
    
    if(var_nodes < replicaNumber) then
      select nodeIP, ipString from Nodes limit 0;
      leave transactionBlock;
   end if;

   -- new replicas
   insert into tempNode
   select nodeIP, ipString
   from Nodes
   where nodeIP <> myIP
   order by replicaCount asc
   limit replicaNumber; 

   -- update counters    
    update Nodes 
   set replicaCount = replicaCount + 1
   where nodeIP in (select repIP from tempNode);
   
    -- update dependencies
    select count(*) from tempNode into n;
    set i=0;
    
    while i<n do
      insert into ReplicaOf(Master, Replica) 
        select myIP, repIP from tempNode limit i,1;
        set i = i+1;
    end while;
    
   -- return nodesIP
   select repIP, repIpStr from tempNode;
                     
   END transactionBlock;

   drop temporary table if exists tempNode;
   commit;
END"""
                     
allOp = [createNodeTable, createReplicaTable, dropIfExistsRep, replicaProc]