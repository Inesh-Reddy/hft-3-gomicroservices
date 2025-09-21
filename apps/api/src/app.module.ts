import { Module } from '@nestjs/common';
import { TickerModule } from './ticker/ticker.module';

@Module({
  imports: [TickerModule],
})
export class AppModule {}
