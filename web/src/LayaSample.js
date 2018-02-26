(function()
{
	var Sprite  = Laya.Sprite;
	var Stage   = Laya.Stage;
    	var Event = Laya.Event;
	var Texture = Laya.Texture;
	var Browser = Laya.Browser;
	var Handler = Laya.Handler;
	var WebGL   = Laya.WebGL;
            var weapon = null;
    var backimg = null;
    var margin = 100;


	(function()
	{
		Laya.init(1560,1020, WebGL);

		Laya.stage.alignV = Stage.ALIGN_MIDDLE;
		Laya.stage.alignH = Stage.ALIGN_CENTER;

		Laya.stage.scaleMode = "showall";
		Laya.stage.bgColor = "#232628";

Laya.loader.load(["buyu/back.jpg","res/atlas/buyu.atlas"],Handler.create(this,loaded));

	})();

	function loaded()
	{
             weapon = new Sprite();
     backimg = new Sprite();

        backimg.loadImage("buyu/back.jpg");
        backimg.size(1560,1020);
                Laya.stage.addChild(backimg);
		weapon.loadImage("buyu/weapon.png",0,0);
        weapon.size(50,50);
        weapon.pos(Laya.stage.width/2-weapon.width,Laya.stage.height-margin-weapon.height);
        weapon.pivot(142,180);
                		Laya.stage.addChild(weapon);
                        Laya.stage.on(Event.MOUSE_DOWN, this, shoot);

	}

    function shoot(e){
        var x = e.stageX-Laya.stage.width/2;
        var y = Laya.stage.height - e.stageY-margin-104/284*50;
        if(y>0){
            angle = 180*Math.atan(x/y)/(Math.PI);
            weapon.rotation = angle;
            var bullet = new Sprite();
            bullet.loadImage("buyu/zidan.png");
            bullet.size(20,20);
            bullet.pos(weapon.x,weapon.y);
            bullet.pivot(28,133);
            Laya.stage.addChild(bullet);
            bullet.rotation = angle;
            Laya.timer.frameLoop(1, this, function(){
                bullet.x +=Math.sin(angle/180*Math.PI)*10;
                bullet.y -=Math.cos(angle/180*Math.PI)*10;
            });
        }
    }

})();