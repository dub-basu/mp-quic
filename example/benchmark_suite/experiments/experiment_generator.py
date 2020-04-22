import os
from shutil import copystat, copyfile

qdisc_start = "---\n"
qdisc_terminate = "...\n"

qdisc_interface = '''
  - 
    interface: "{}"
    
    netem:
      netem_on: yes
      
      packet_corruption_on: no
      packet_corruption: 1%

      delay_on: yes
      delay: 
        delay_time: {}ms
        error: {}ms
        correlation: 25%

      packet_loss_on: yes
      packet_loss: 
        loss: {}%
        correlation: 25%
   
    tbf: 
      tbf_on: yes
      rate: {}kbit
      latency: '3000'
      burst: 32kbit
#---------------------------#
  
'''

exp = '''
# Blocks of size BLOCK_SIZE bytes will be requested 
# DELAY_TIME milliseconds pause after each block
# ITERATIONS number of times, with 

BLOCK_SIZE: {}
DELAY_TIME: {}
ITERATIONS: {}

varying_qdisc: no
'''




BLOCK_SIZES = [64,128,192,256,512,768,1024,2048,2560,3072,3574,4096,5120]
TIME_REQ = [200,400,600,800]
ITERATIONS = [10,50]

# s1_delays = [50]
# s1_losses = [1]
# s1_bandwidths = [2000] #2Mbps

# s2_delays = [500]
# s2_losses = [0.001]
# s2_bandwidths = [250] #250kbps


s1_delays = [50]
s1_losses = [1]
s1_bandwidths = [2000] #2Mbps

s2_delays = [50]
s2_losses = [1]
s2_bandwidths = [2000] #250kbps
filename_begin = "{}"

exp_num = 0

for s1_delay in s1_delays:
    for s1_loss in s1_losses:
        for s1_bandwidth in s1_bandwidths:
            for s2_delay in s2_delays:
                for s2_loss in s2_losses:
                    for s2_bandwidth in s2_bandwidths:
                        for block_size in BLOCK_SIZES:
                            for time_req in TIME_REQ:
                                for iterations in ITERATIONS:
                                    exp_num += 1
                                    exp_dir = filename_begin.format(str(exp_num)) 
                                    os.mkdir( exp_dir )
                                    os.mkdir( exp_dir + '/client')
                                    os.mkdir( exp_dir + '/router')
                                    os.mkdir( exp_dir + '/server')

                                    copyfile('experiment_demo/topo.py', exp_dir+'/topo.py')
                                    copystat('experiment_demo/topo.py', exp_dir+'/topo.py')
                                    copyfile('experiment_demo/client/config.sh', exp_dir+'/client/config.sh')
                                    copystat('experiment_demo/client/config.sh', exp_dir+'/client/config.sh')
                                    copyfile('experiment_demo/server/config.sh', exp_dir+'/server/config.sh')
                                    copystat('experiment_demo/server/config.sh', exp_dir+'/server/config.sh')
                                    copyfile('experiment_demo/router/config.sh', exp_dir+'/router/config.sh')
                                    copystat('experiment_demo/router/config.sh', exp_dir+'/router/config.sh')

                                    qdisc_file = open(exp_dir+"/qdisc.yaml", "w+")
                                    qdisc_file.write(qdisc_start)
                                    
                                    qdisc_file.write( qdisc_interface.format("s1-eth1", str(s1_delay), str(s1_delay//10), str(s1_loss), str(s1_bandwidth)) )
                                    qdisc_file.write( qdisc_interface.format("s1-eth2", str(s1_delay), str(s1_delay//10), str(s1_loss), str(s1_bandwidth)) )
                                    qdisc_file.write( qdisc_interface.format("s2-eth1", str(s2_delay), str(s2_delay//10), str(s2_loss), str(s2_bandwidth)) )
                                    qdisc_file.write( qdisc_interface.format("s2-eth2", str(s2_delay), str(s2_delay//10), str(s2_loss), str(s2_bandwidth)) )

                                    qdisc_file.write(qdisc_terminate)
                                    qdisc_file.close()

                                    
                                    exp_file = open(exp_dir+"/exp.yaml", "w+")
                                    exp_file.write(exp.format(str(block_size),str(time_req), str(iterations)))
                                    exp_file.close()