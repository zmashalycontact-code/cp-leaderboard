import React, { useMemo } from 'react';
import { motion } from 'framer-motion';
import { 
  BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell, 
  PieChart, Pie, ScatterChart, Scatter, ZAxis, Legend 
} from 'recharts';
import { Target, Users, Zap } from 'lucide-react';

const RANK_COLORS = {
  'CP MASTER': '#eab308', // Yellow
  'The GOAT': '#d946ef',  // Fuchsia
  'تنين مجنح': '#ef4444',   // Red
  'شكسبير': '#a855f7',    // Purple
  'باشا ستراكشر': '#14b8a6', // Teal
  'روش': '#22c55e',       // Green
  'كحيان': '#94a3b8'      // Slate
};

const CustomTooltip = ({ active, payload }) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;
    return (
      <div className="bg-slate-900 border border-slate-700 p-4 rounded-xl shadow-[0_0_20px_rgba(0,0,0,0.8)]">
        <p className="text-white font-black mb-1 border-b border-white/10 pb-2">{data.name || data.handle}</p>
        {payload.map((entry, index) => (
          <p key={index} style={{ color: entry.color || entry.fill }} className="text-sm font-bold mt-1">
            {entry.name === 'value' ? 'العدد' : entry.name}: {entry.value}
          </p>
        ))}
      </div>
    );
  }
  return null;
};

export default function StatsDashboard({ users }) {
  // 1. حساب الإحصائيات العامة (Summary Cards)
  const totalCommunityPoints = users.reduce((acc, curr) => acc + (curr.season_points || 0), 0).toFixed(0);
  const totalCommunityHidden = users.reduce((acc, curr) => acc + (curr.hidden_solved || 0), 0);
  const totalActiveUsers = users.filter(u => u.activity_7d > 0).length;

  // 2. توزيع الرتب للمجتمع (Donut Chart)
  const rankDistribution = useMemo(() => {
    const counts = users.reduce((acc, user) => {
      const tier = user.rank_tier || 'كحيان';
      acc[tier] = (acc[tier] || 0) + 1;
      return acc;
    }, {});
    return Object.keys(counts).map(key => ({
      name: key,
      value: counts[key]
    })).sort((a, b) => b.value - a.value);
  }, [users]);

  // 3. أبطال المعافرة - Horizontal Bar (أكثر ناس جابت AC بعد محاولات كتير)
  const topStrugglers = useMemo(() => {
    return [...users]
      .filter(u => u.struggle_count > 0)
      .sort((a, b) => b.struggle_count - a.struggle_count)
      .slice(0, 8)
      .map(u => ({ handle: u.handle, 'بطلوع الروح': u.struggle_count }));
  }, [users]);

  // 4. العلاقة بين الريت والنقاط - Scatter Plot (Deep Insight)
  const ratingVsPoints = useMemo(() => {
    return users
      .filter(u => u.current_rating > 0) // بناخد اللي ليهم ريت بس
      .map(u => ({
        handle: u.handle,
        rating: u.current_rating,
        points: parseFloat((u.season_points || 0).toFixed(1)),
        tier: u.rank_tier
      }));
  }, [users]);

  return (
    <motion.div 
      initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      {/* 🃏 كروت الملخص السريعة */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-gradient-to-br from-indigo-900/50 to-slate-900 border border-indigo-500/30 p-6 rounded-3xl flex items-center gap-4">
          <div className="p-4 bg-indigo-500/20 rounded-2xl text-indigo-400"><Target size={32} /></div>
          <div>
            <p className="text-slate-400 text-sm font-bold">إجمالي نقاط المعسكر</p>
            <p className="text-3xl font-black text-white">{totalCommunityPoints} <span className="text-sm text-indigo-400">PTS</span></p>
          </div>
        </div>
        <div className="bg-gradient-to-br from-emerald-900/50 to-slate-900 border border-emerald-500/30 p-6 rounded-3xl flex items-center gap-4">
          <div className="p-4 bg-emerald-500/20 rounded-2xl text-emerald-400"><Users size={32} /></div>
          <div>
            <p className="text-slate-400 text-sm font-bold">المنافسين النشطين (آخر أسبوع)</p>
            <p className="text-3xl font-black text-white">{totalActiveUsers} <span className="text-sm text-emerald-400">CPers</span></p>
          </div>
        </div>
        <div className="bg-gradient-to-br from-purple-900/50 to-slate-900 border border-purple-500/30 p-6 rounded-3xl flex items-center gap-4">
          <div className="p-4 bg-purple-500/20 rounded-2xl text-purple-400"><Zap size={32} /></div>
          <div>
            <p className="text-slate-400 text-sm font-bold">شيتات تم تدميرها (مخفية)</p>
            <p className="text-3xl font-black text-white">{totalCommunityHidden} <span className="text-sm text-purple-400">مسألة</span></p>
          </div>
        </div>
      </div>

      {/* 📊 شبكة الرسومات العميقة */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        
        {/* الرسمة الأولى: التوزيع الطبقي للمجتمع (Donut Chart) */}
        <div className="bg-white/[0.02] border border-white/10 p-6 rounded-3xl">
          <h2 className="text-lg font-black text-slate-300 mb-4 flex items-center gap-2">
            👑 التوزيع الطبقي للمعسكر
          </h2>
          <div className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={rankDistribution}
                  cx="50%" cy="50%"
                  innerRadius={70} outerRadius={110}
                  paddingAngle={5}
                  dataKey="value"
                  stroke="none"
                >
                  {rankDistribution.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={RANK_COLORS[entry.name] || '#ffffff'} />
                  ))}
                </Pie>
                <Tooltip content={<CustomTooltip />} />
                <Legend verticalAlign="bottom" height={36} wrapperStyle={{ fontSize: '12px', fontWeight: 'bold' }}/>
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* الرسمة التانية: أبطال المعافرة (Horizontal Bar) - مرنة لعدد كبير */}
        <div className="bg-white/[0.02] border border-white/10 p-6 rounded-3xl">
          <h2 className="text-lg font-black text-slate-300 mb-4 flex items-center gap-2">
            💀 أبطال المعافرة (أكثر ناس بتحاول وتغلط لحد ما تجيب AC)
          </h2>
          <div className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={topStrugglers} layout="vertical" margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
                <XAxis type="number" tick={{ fill: '#64748b', fontSize: 12 }} />
                <YAxis dataKey="handle" type="category" tick={{ fill: '#cbd5e1', fontSize: 11, fontWeight: 'bold' }} width={90} />
                <Tooltip content={<CustomTooltip />} cursor={{ fill: 'rgba(255,255,255,0.05)' }} />
                <Bar dataKey="بطلوع الروح" radius={[0, 4, 4, 0]} barSize={20}>
                  {topStrugglers.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill="#f87171" />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* الرسمة التالتة: الذكاء vs المجهود (Scatter Plot) - بتاخد العرض كله */}
        <div className="bg-white/[0.02] border border-white/10 p-6 rounded-3xl lg:col-span-2">
          <h2 className="text-lg font-black text-slate-300 mb-4 flex items-center gap-2">
            📈 المجهود مقابل التقييم (هل الريت العالي بيجيب نقط أكتر؟)
          </h2>
          <p className="text-xs text-slate-500 mb-6 font-bold">كل نقطة بتمثل متسابق. المحور الأفقي هو التقييم (Rating) والرأسي هو النقاط في السيزون.</p>
          <div className="h-[350px]">
            <ResponsiveContainer width="100%" height="100%">
              <ScatterChart margin={{ top: 20, right: 20, bottom: 20, left: -10 }}>
                <XAxis type="number" dataKey="rating" name="CF Rating" domain={['dataMin - 100', 'dataMax + 100']} tick={{ fill: '#64748b' }} stroke="#334155" />
                <YAxis type="number" dataKey="points" name="النقاط" tick={{ fill: '#64748b' }} stroke="#334155" />
                <ZAxis type="number" range={[100, 100]} />
                <Tooltip cursor={{ strokeDasharray: '3 3', stroke: '#475569' }} content={<CustomTooltip />} />
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
