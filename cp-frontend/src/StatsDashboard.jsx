import React, { useMemo } from 'react';
import { motion } from 'framer-motion';
import { 
  ComposedChart, Line, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell, 
  ScatterChart, Scatter, ZAxis, CartesianGrid
} from 'recharts';
import { Target, Users, Zap, Flame, TrendingUp, BarChart3, Swords } from 'lucide-react';

const RANK_COLORS = {
  'CP MASTER': '#eab308', 'The GOAT': '#d946ef', 'تنين مجنح': '#ef4444', 
  'شكسبير': '#a855f7', 'باشا ستراكشر': '#14b8a6', 'روش': '#22c55e', 'كحيان': '#94a3b8'
};

const CustomTooltip = ({ active, payload, label }) => {
  if (active && payload && payload.length) {

    const title = payload[0].payload.handle || label;

    return (
      <div className="bg-slate-900 border border-slate-700 p-4 rounded-xl shadow-[0_0_30px_rgba(0,0,0,0.9)] z-50 relative" dir="rtl">
        <p className="text-white font-black mb-2 border-b border-white/10 pb-2 text-center text-lg">
          {title}
        </p>
        <div className="space-y-2">
          {payload.map((entry, index) => (
            <div key={index} className="flex justify-between items-center gap-8">
               <span className="text-sm font-bold" style={{ color: entry.color || entry.fill }}>
                {entry.name}:
              </span>
              <span className="text-white font-black">
                {/* Fallback لضمان ظهور الرقم الصحيح في كل أنواع الرسومات */}
                {entry.payload[entry.dataKey + "_raw"] ?? entry.value} 
              </span>
            </div>
          ))}
        </div>
      </div>
    );
  }
  return null;
};

// كومبوننت كارت الإحصائيات السريعة
const StatCard = ({ icon, label, val, color }) => (
  <div className={`bg-gradient-to-br from-${color}-900/50 to-slate-900 border border-${color}-500/30 p-6 rounded-[2rem] flex items-center gap-4 shadow-lg`}>
    <div className={`p-4 bg-${color}-500/20 rounded-2xl text-${color}-400`}>{icon}</div>
    <div>
      <p className="text-slate-400 text-[11px] font-bold uppercase tracking-wider">{label}</p>
      <p className="text-3xl font-black text-white">{val}</p>
    </div>
  </div>
);

export default function StatsDashboard({ users }) {
  if (!users || users.length === 0) return null;

  // 1. بيانات الزخم التاريخي (كما هي بطلبك)
  const momentumData = useMemo(() => [...users].sort((a, b) => b.total_solved - a.total_solved), [users]);
  
  // 2. بيانات المجهود مقابل التقييم (كما هي بطلبك)
  const ratingVsPoints = useMemo(() => users.filter(u => u.current_rating > 0).map(u => ({
    handle: u.handle, rating: u.current_rating, points: parseFloat((u.season_points || 0).toFixed(1)), tier: u.rank_tier
  })), [users]);

  // 3. الرسمة الجديدة العبقرية (ميزان الحصاد والمعافرة) ببيانات حقيقية 100%
  const harvestData = useMemo(() => {
    return [...users]
      .sort((a, b) => b.season_points - a.season_points) // ترتيب بأعلى النقاط
      .slice(0, 7) // عرض أعلى 7 متسابقين
      .map(u => ({
        handle: u.handle,
        'النقاط': parseFloat((u.season_points || 0).toFixed(1)),
        'المعافرة': u.struggle_count || 0
      }));
  }, [users]);

  return (
    <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-10 pb-20" dir="rtl">
      
      {/* كروت الإحصائيات السريعة */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <StatCard icon={<Target size={24}/>} label="نقاط السيزون" val={users.reduce((acc,u)=>acc+(u.season_points||0),0).toFixed(0)} color="indigo" />
        <StatCard icon={<Flame size={24}/>} label="نشاط المعسكر" val={users.reduce((acc,u)=>acc+(u.activity_7d||0),0)} color="orange" />
        <StatCard icon={<Users size={24}/>} label="النشطين حالياً" val={users.filter(u=>u.activity_7d>0).length} color="emerald" />
        <StatCard icon={<Zap size={24}/>} label="شيتات مخفية" val={users.reduce((acc,u)=>acc+(u.hidden_solved||0),0)} color="purple" />
      </div>

      {/* 🔥 الرسمة الجديدة العبقرية: ميزان الحصاد والمعافرة */}
      <div className="bg-white/[0.02] border border-white/10 p-10 rounded-[2.5rem] shadow-2xl relative overflow-hidden">
        <div className="absolute top-0 right-0 w-64 h-64 bg-rose-500/5 blur-[100px] rounded-full pointer-events-none" />
        
        <div className="text-center mb-10">
          <h2 className="text-3xl font-black text-white mb-2 flex items-center justify-center gap-3">
            <Swords size={32} className="text-rose-400" /> ميزان الحصاد والمعافرة
          </h2>
          <p className="text-slate-500 font-bold">مقارنة دقيقة لأعلى المتصدرين: حجم النقاط المحصودة مقابل كمية المعافرة (WA/TLE) للوصول إليها.</p>
        </div>

        <div className="h-[450px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <ComposedChart data={harvestData} margin={{ top: 20, right: 20, left: 20, bottom: 20 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} opacity={0.4} />
              <XAxis 
                dataKey="handle" 
                tick={{ fill: '#94a3b8', fontSize: 12, fontWeight: 'bold' }} 
                tickMargin={15} 
                axisLine={false}
                tickLine={false}
              />
              <YAxis 
                yAxisId="left" 
                tick={{ fill: '#eab308', fontSize: 12, fontWeight: 'bold' }} 
                tickMargin={15} 
                axisLine={false} 
                tickLine={false} 
              />
              <YAxis 
                yAxisId="right" 
                orientation="right" 
                tick={{ fill: '#f43f5e', fontSize: 12, fontWeight: 'bold' }} 
                tickMargin={15} 
                axisLine={false} 
                tickLine={false} 
              />
              <Tooltip content={<CustomTooltip />} cursor={{ fill: 'rgba(255,255,255,0.03)' }} />
              
              <Bar yAxisId="left" dataKey="النقاط" name="النقاط المحصودة" fill="#eab308" radius={[8, 8, 0, 0]} barSize={40} />
              <Line yAxisId="right" type="monotone" dataKey="المعافرة" name="مرات المعافرة" stroke="#f43f5e" strokeWidth={4} dot={{ r: 6, fill: '#f43f5e', stroke: '#fff', strokeWidth: 2 }} />
            </ComposedChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* 📊 الرسمتين اللي قولت ملمسهمش (الزخم والمجهود) */}
      <div className="space-y-10">
         <div className="bg-white/[0.02] border border-white/10 p-6 rounded-3xl">
          <h2 className="text-lg font-black text-slate-300 mb-2 flex items-center gap-2"><TrendingUp size={20} className="text-blue-400" /> الزخم التاريخي مقابل النشاط الحالي</h2>
          <div className="h-[450px]">
            <ResponsiveContainer width="100%" height="100%">
              <ComposedChart data={momentumData} margin={{ top: 20, right: 20, left: 40, bottom: 20 }}>
                <XAxis dataKey="handle" tick={{ fill: '#94a3b8', fontSize: 10, fontWeight: 'bold', dx: -15, dy: 55 }} tickFormatter={(val) => val.substring(0, 12)} angle={-45} textAnchor="end" interval={0} height={80} />
                <YAxis yAxisId="left" tick={{ fill: '#94a3b8', fontSize: 12, fontWeight: 'bold' }} tickMargin={20} width={65} />
                <YAxis yAxisId="right" orientation="right" tick={{ fill: '#10b981', fontSize: 12, fontWeight: 'bold' }} tickMargin={15} />
                <Tooltip content={<CustomTooltip />} />
                <Bar yAxisId="left" dataKey="total_solved" name="تاريخ الحل" fill="#3b82f6" radius={[4, 4, 0, 0]} />
                <Line yAxisId="right" type="monotone" dataKey="activity_7d" name="نشاط الأسبوع" stroke="#10b981" strokeWidth={4} dot={{ r: 4, fill: '#10b981', stroke: '#fff' }} />
              </ComposedChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white/[0.02] border border-white/10 p-6 rounded-3xl">
          <h2 className="text-lg font-black text-slate-300 mb-2 flex items-center gap-2"><BarChart3 size={20} className="text-orange-400" /> المجهود مقابل التقييم</h2>
          <div className="h-[450px]">
            <ResponsiveContainer width="100%" height="100%">
              <ScatterChart margin={{ top: 20, right: 30, bottom: 20, left: -10 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" opacity={0.5} />
                <XAxis type="number" dataKey="rating" domain={['dataMin', 'dataMax']} tick={{ fill: '#94a3b8', fontSize: 12, fontWeight: 'bold' }} stroke="#334155" tickMargin={15} />
                <YAxis type="number" dataKey="points" domain={[0, 'auto']} tick={{ fill: '#94a3b8', fontSize: 12, fontWeight: 'bold' }} stroke="#334155" tickMargin={10} />
                <ZAxis type="number" range={[100, 100]} />
                <Tooltip cursor={{ strokeDasharray: '3 3' }} content={<CustomTooltip />} />
                <Scatter name="المتسابقين" data={ratingVsPoints}>
                  {ratingVsPoints.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={RANK_COLORS[entry.tier] || '#60a5fa'} />
                  ))}
                </Scatter>
              </ScatterChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

    </motion.div>
  );
}
