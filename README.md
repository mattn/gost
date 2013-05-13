# gost

github push deployer 

                           ;####+.                    
                        `##########.                  
                       '####';:;####+                 
                      ###'........'###                
                     ###,..........:###               
                    +##.............;##+              
                   .##:..............+##`             
                   ###.`..............##+             
                  .##.```.........```.###             
                  ###.```..,.,.....``.;##,            
                  ##+````.;',,.;,..``..##+            
                 ,##:````;##',:##:.```.###            
                 +##.``..###',;###.```.###            
                 ###.`..####;,;###;```.###            
                 ###...####+,,,+###..`.###`           
                 ##+`.#####;,.,'####,..+##.           
                `##'.;####',...,'####,.+##.           
                `##'.,###;:..:..,+###,.+##,           
                .##+..,:,,,.+';.:,,,...###,           
                `###.,,,,:.,#;#..,,,,,.###`           
                 ###',,,,...,,,.`,,,,,;###            
                 ####:,,..........,,:,####            
                 :####.....,++;....,,;###;            
                 '####,....####'...,,####             
                 #####+...:#####,...:####             
                 ######,..:#####,..,#####             
                 ######;..,#####,..,#####             
                 #######..,#####,.,+#####.            
                 #######,..#####..,######;            
                `#######'..#####..;#######            
                :########..'###'.,########            
                '########,.;###;.,########            
                #########,.,###,,:########            
                #########'..###.,+########,           
                #+########,.###.,#########;           
                ##########,.'#',:##########           
                ##########:,,+,,'##########           
               ,##+++######:,,::########+##           
               +###++#######';'#########+##           
              `++##########################:          
              ++############################          
             +##############################,         
           ,+#############+##################.        
          :+############++#####+##############:       
          `+++++#####+`+++###+##+.;##'########+       
            ;+++#+#'                    +#####;       
             `'+#                        .+#.         
                ,                                     

## usage

    $ gost -c config.json

## setting

    
    {
        "addr": "127.0.0.1:8889",
        "rpc": "127.0.0.1:5555",
        "apps": {
            "my-app1": {
                "proc": "today_bot",
                "path": "/home/mattn/dev/today_bot/"
            },
            "my-app2": {
                "proc": "today_bot",
                "path": "/home/mattn/dev/today_bot/"
            }
        }
    }

* addr: server address
* rpc: goreman port
* my-app1: this mean github repository name
* proc: goreman proc name to start or stop
* path: path that is working the proc

## tips

You can start gost with goreman/gorem like follow:

goreman configuration(Procfile):

    gorem: /home/mattn/dev/gorem/gorem -c /home/mattn/dev/gorem/config.json
    gost: /home/mattn/dev/gost/gost -c /home/mattn/dev/gost/config.json
    my-app1: ruby /home/mattn/dev/my-app1/web.rb

gorem configuration(config.json):

    {
        "mattn": {
            "address": "127.0.0.1:80",
            "entries": [
                { "path": "/my-app1/", "backend": "http://localhost:8888" },
                { "path": "/deploy/", "backend": "http://localhost:"8889" },
                { "path": "/",    "backend": "/home/mattn/blog/_site" }
            ]
        }
    }

Then, the things that you should do is:

    $ goreman start

Set webhook settings on github settings to:

    http://example.com/deply/

